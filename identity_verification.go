package synnergy

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"sort"
	"sync"
	"time"
)

// IdentityInfo contains basic identity metadata.
type IdentityInfo struct {
	Name        string
	DateOfBirth string
	Nationality string
}

// IdentityStatus represents the lifecycle state of an identity record.
type IdentityStatus string

const (
	// IdentityStatusPending indicates the identity was registered but has not been verified yet.
	IdentityStatusPending IdentityStatus = "pending"
	// IdentityStatusVerified marks the identity as having passed verification.
	IdentityStatusVerified IdentityStatus = "verified"
	// IdentityStatusRevoked marks identities that have been revoked or suspended.
	IdentityStatusRevoked IdentityStatus = "revoked"
)

// VerificationLog records a verification attempt.
type VerificationLog struct {
	Method       string
	Timestamp    time.Time
	EvidenceHash string
	Signer       string
	Notes        string
	Successful   bool
	Metadata     map[string]string
}

// IdentityRecord stores the authoritative state for a registered identity.
type IdentityRecord struct {
	Info         IdentityInfo
	Status       IdentityStatus
	RegisteredAt time.Time
	UpdatedAt    time.Time
	Revision     uint64
	Logs         []VerificationLog
	Metadata     map[string]string
}

// IdentityEventType enumerates identity lifecycle notifications.
type IdentityEventType string

const (
	// IdentityEventRegistered is emitted when a new identity is registered.
	IdentityEventRegistered IdentityEventType = "registered"
	// IdentityEventVerified is emitted when an identity passes verification.
	IdentityEventVerified IdentityEventType = "verified"
	// IdentityEventRevoked is emitted when an identity is revoked.
	IdentityEventRevoked IdentityEventType = "revoked"
)

// IdentityEvent is dispatched to watchers when an identity changes.
type IdentityEvent struct {
	Type    IdentityEventType
	Address string
	Record  IdentityRecord
}

// IdentityEventHandler receives identity lifecycle events.
type IdentityEventHandler interface {
	HandleIdentityEvent(ctx context.Context, event IdentityEvent) error
}

// IdentityEventHandlerFunc adapts a function into an IdentityEventHandler.
type IdentityEventHandlerFunc func(context.Context, IdentityEvent) error

// HandleIdentityEvent implements IdentityEventHandler.
func (f IdentityEventHandlerFunc) HandleIdentityEvent(ctx context.Context, event IdentityEvent) error {
	return f(ctx, event)
}

// IdentityServiceOption configures IdentityService behaviour.
type IdentityServiceOption func(*IdentityService)

// WithIdentityEventHandler registers a handler for identity lifecycle events.
func WithIdentityEventHandler(handler IdentityEventHandler) IdentityServiceOption {
	return func(s *IdentityService) {
		if handler != nil {
			s.watchers = append(s.watchers, handler)
		}
	}
}

// WithIdentityLogRetention configures how many verification entries are retained per identity.
func WithIdentityLogRetention(limit int) IdentityServiceOption {
	return func(s *IdentityService) {
		if limit > 0 {
			s.maxLogs = limit
		}
	}
}

var (
	// ErrAddressRequired is returned when an operation is attempted with an empty address.
	ErrAddressRequired = errors.New("identity: address required")
	// ErrIdentityExists indicates the identity is already registered.
	ErrIdentityExists = errors.New("identity: already registered")
	// ErrIdentityNotFound indicates the requested identity does not exist.
	ErrIdentityNotFound = errors.New("identity: not registered")
	// ErrMethodRequired indicates that a verification method must be supplied.
	ErrMethodRequired = errors.New("identity: verification method required")
	// ErrInvalidSignature indicates a provided verification signature failed validation.
	ErrInvalidSignature = errors.New("identity: invalid verification signature")
)

// IdentityService manages verified addresses on the ledger.
type IdentityService struct {
	mu       sync.RWMutex
	records  map[string]IdentityRecord
	watchers []IdentityEventHandler
	maxLogs  int
}

// NewIdentityService creates a new IdentityService instance.
func NewIdentityService(opts ...IdentityServiceOption) *IdentityService {
	svc := &IdentityService{
		records: make(map[string]IdentityRecord),
		maxLogs: 64,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(svc)
		}
	}
	return svc
}

// Register stores identity information for an address.
func (s *IdentityService) Register(addr, name, dob, nationality string) error {
	return s.RegisterContext(context.Background(), addr, IdentityInfo{Name: name, DateOfBirth: dob, Nationality: nationality}, nil)
}

// RegisterContext stores identity information with optional metadata using the provided context.
func (s *IdentityService) RegisterContext(ctx context.Context, addr string, info IdentityInfo, metadata map[string]string) error {
	if addr == "" {
		return ErrAddressRequired
	}
	now := time.Now().UTC()
	record := IdentityRecord{
		Info:         info,
		Status:       IdentityStatusPending,
		RegisteredAt: now,
		UpdatedAt:    now,
		Revision:     1,
		Metadata:     cloneMetadata(metadata),
	}

	s.mu.Lock()
	if _, exists := s.records[addr]; exists {
		s.mu.Unlock()
		return ErrIdentityExists
	}
	s.records[addr] = record
	s.mu.Unlock()

	return s.fireEvent(ctx, IdentityEvent{Type: IdentityEventRegistered, Address: addr, Record: copyIdentityRecord(record)})
}

// VerificationRequest contains advanced verification parameters used by VerifyContext.
type VerificationRequest struct {
	Method           string
	EvidenceHash     []byte
	Notes            string
	SignerPublicKey  ed25519.PublicKey
	Signature        []byte
	Timestamp        time.Time
	Metadata         map[string]string
	TransitionStatus IdentityStatus
}

// Verify records a verification method for an address.
func (s *IdentityService) Verify(addr, method string) error {
	req := VerificationRequest{Method: method}
	return s.VerifyContext(context.Background(), addr, req)
}

// VerifyContext records a verification method and optional signature/evidence metadata.
func (s *IdentityService) VerifyContext(ctx context.Context, addr string, req VerificationRequest) error {
	if addr == "" {
		return ErrAddressRequired
	}
	if req.Method == "" {
		return ErrMethodRequired
	}
	ts := req.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	if len(req.Signature) > 0 || len(req.SignerPublicKey) > 0 {
		if len(req.Signature) != ed25519.SignatureSize || len(req.SignerPublicKey) != ed25519.PublicKeySize {
			return ErrInvalidSignature
		}
		digest := verificationDigest(addr, req)
		if !ed25519.Verify(req.SignerPublicKey, digest, req.Signature) {
			return ErrInvalidSignature
		}
	}

	s.mu.Lock()
	record, ok := s.records[addr]
	if !ok {
		s.mu.Unlock()
		return ErrIdentityNotFound
	}
	record.Status = desiredStatus(record.Status, req.TransitionStatus)
	record.UpdatedAt = ts
	record.Revision++
	if len(req.Metadata) > 0 {
		if record.Metadata == nil {
			record.Metadata = make(map[string]string, len(req.Metadata))
		}
		for k, v := range req.Metadata {
			record.Metadata[k] = v
		}
	}
	logEntry := VerificationLog{
		Method:       req.Method,
		Timestamp:    ts,
		EvidenceHash: encodeEvidence(req.EvidenceHash),
		Signer:       encodeSigner(req.SignerPublicKey),
		Notes:        req.Notes,
		Successful:   true,
		Metadata:     cloneMetadata(req.Metadata),
	}
	record.Logs = append(record.Logs, logEntry)
	if s.maxLogs > 0 && len(record.Logs) > s.maxLogs {
		record.Logs = record.Logs[len(record.Logs)-s.maxLogs:]
	}
	s.records[addr] = record
	s.mu.Unlock()

	return s.fireEvent(ctx, IdentityEvent{Type: IdentityEventVerified, Address: addr, Record: copyIdentityRecord(record)})
}

// Revoke marks an identity as revoked while recording the supplied reason.
func (s *IdentityService) Revoke(ctx context.Context, addr, reason string, metadata map[string]string) error {
	if addr == "" {
		return ErrAddressRequired
	}
	now := time.Now().UTC()

	s.mu.Lock()
	record, ok := s.records[addr]
	if !ok {
		s.mu.Unlock()
		return ErrIdentityNotFound
	}
	record.Status = IdentityStatusRevoked
	record.UpdatedAt = now
	record.Revision++
	if metadata != nil {
		if record.Metadata == nil {
			record.Metadata = make(map[string]string, len(metadata))
		}
		for k, v := range metadata {
			record.Metadata[k] = v
		}
	}
	record.Logs = append(record.Logs, VerificationLog{
		Method:     "revoke",
		Timestamp:  now,
		Notes:      reason,
		Successful: false,
		Metadata:   cloneMetadata(metadata),
	})
	if s.maxLogs > 0 && len(record.Logs) > s.maxLogs {
		record.Logs = record.Logs[len(record.Logs)-s.maxLogs:]
	}
	s.records[addr] = record
	s.mu.Unlock()

	return s.fireEvent(ctx, IdentityEvent{Type: IdentityEventRevoked, Address: addr, Record: copyIdentityRecord(record)})
}

// Info retrieves identity information for an address.
func (s *IdentityService) Info(addr string) (IdentityInfo, bool) {
	s.mu.RLock()
	record, ok := s.records[addr]
	s.mu.RUnlock()
	if !ok {
		return IdentityInfo{}, false
	}
	return record.Info, true
}

// Status returns the status of the identity if present.
func (s *IdentityService) Status(addr string) (IdentityStatus, bool) {
	s.mu.RLock()
	record, ok := s.records[addr]
	s.mu.RUnlock()
	if !ok {
		return "", false
	}
	return record.Status, true
}

// Logs returns verification logs for an address.
func (s *IdentityService) Logs(addr string) []VerificationLog {
	s.mu.RLock()
	record, ok := s.records[addr]
	s.mu.RUnlock()
	if !ok {
		return nil
	}
	return copyLogs(record.Logs)
}

// Record returns a defensive copy of an identity record if it exists.
func (s *IdentityService) Record(addr string) (IdentityRecord, bool) {
	s.mu.RLock()
	record, ok := s.records[addr]
	s.mu.RUnlock()
	if !ok {
		return IdentityRecord{}, false
	}
	return copyIdentityRecord(record), true
}

// Addresses returns a snapshot of all registered addresses.
func (s *IdentityService) Addresses() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.records))
	for addr := range s.records {
		out = append(out, addr)
	}
	sort.Strings(out)
	return out
}

func (s *IdentityService) fireEvent(ctx context.Context, event IdentityEvent) error {
	s.mu.RLock()
	watchers := append([]IdentityEventHandler(nil), s.watchers...)
	s.mu.RUnlock()
	var errs []error
	for _, watcher := range watchers {
		if watcher == nil {
			continue
		}
		if err := watcher.HandleIdentityEvent(ctx, event); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func cloneMetadata(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func copyLogs(in []VerificationLog) []VerificationLog {
	if len(in) == 0 {
		return nil
	}
	out := make([]VerificationLog, len(in))
	for i, log := range in {
		out[i] = VerificationLog{
			Method:       log.Method,
			Timestamp:    log.Timestamp,
			EvidenceHash: log.EvidenceHash,
			Signer:       log.Signer,
			Notes:        log.Notes,
			Successful:   log.Successful,
			Metadata:     cloneMetadata(log.Metadata),
		}
	}
	return out
}

func copyIdentityRecord(in IdentityRecord) IdentityRecord {
	return IdentityRecord{
		Info:         in.Info,
		Status:       in.Status,
		RegisteredAt: in.RegisteredAt,
		UpdatedAt:    in.UpdatedAt,
		Revision:     in.Revision,
		Logs:         copyLogs(in.Logs),
		Metadata:     cloneMetadata(in.Metadata),
	}
}

func verificationDigest(addr string, req VerificationRequest) []byte {
	h := sha256.New()
	h.Write([]byte(addr))
	h.Write([]byte("|"))
	h.Write([]byte(req.Method))
	h.Write([]byte("|"))
	h.Write(req.EvidenceHash)
	keys := make([]string, 0, len(req.Metadata))
	for k := range req.Metadata {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		h.Write([]byte("|"))
		h.Write([]byte(k))
		h.Write([]byte("="))
		h.Write([]byte(req.Metadata[k]))
	}
	if req.Timestamp.IsZero() {
		h.Write([]byte("|ts:"))
	} else {
		h.Write([]byte("|" + req.Timestamp.UTC().Format(time.RFC3339Nano)))
	}
	return h.Sum(nil)
}

func desiredStatus(current IdentityStatus, requested IdentityStatus) IdentityStatus {
	switch requested {
	case IdentityStatusPending, IdentityStatusVerified, IdentityStatusRevoked:
		return requested
	case "":
		switch current {
		case "":
			return IdentityStatusPending
		case IdentityStatusPending:
			return IdentityStatusVerified
		default:
			return current
		}
	default:
		return current
	}
}

func encodeEvidence(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return base64.RawStdEncoding.EncodeToString(b)
}

func encodeSigner(pub ed25519.PublicKey) string {
	if len(pub) == 0 {
		return ""
	}
	return base64.RawStdEncoding.EncodeToString(pub)
}
