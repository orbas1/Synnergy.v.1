package governance

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// AuditEvent captures the context for a governance action.
type AuditEvent struct {
	ID        string            `json:"id"`
	Actor     string            `json:"actor"`
	Action    string            `json:"action"`
	Scope     string            `json:"scope"`
	NodeID    string            `json:"node_id"`
	Reason    string            `json:"reason"`
	GasBudget uint64            `json:"gas_budget"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// AuditRecord represents an immutable entry in the governance log.
type AuditRecord struct {
	Sequence   uint64 `json:"sequence"`
	PrevHash   string `json:"prev_hash"`
	Hash       string `json:"hash"`
	Signature  string `json:"signature,omitempty"`
	AuditEvent `json:",inline"`
}

// Signer generates and validates signatures for audit entries.
type Signer interface {
	Sign(message []byte) (string, error)
	Verify(message []byte, signature string) bool
}

// HMACSigner implements Signer using HMAC-SHA256.
type HMACSigner struct {
	key []byte
}

// NewHMACSigner instantiates a signer from the supplied secret.
func NewHMACSigner(key []byte) *HMACSigner {
	buf := make([]byte, len(key))
	copy(buf, key)
	return &HMACSigner{key: buf}
}

// Sign produces a base64 encoded HMAC signature.
func (s *HMACSigner) Sign(message []byte) (string, error) {
	mac := hmac.New(sha256.New, s.key)
	if _, err := mac.Write(message); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// Verify validates a base64 encoded HMAC signature.
func (s *HMACSigner) Verify(message []byte, signature string) bool {
	mac := hmac.New(sha256.New, s.key)
	if _, err := mac.Write(message); err != nil {
		return false
	}
	expected := mac.Sum(nil)
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	return hmac.Equal(expected, sig)
}

// Observer receives notifications for newly appended audit records.
type Observer func(AuditRecord)

// AuditLog stores tamper evident governance events with optional signing.
type AuditLog struct {
	mu          sync.RWMutex
	entries     []AuditRecord
	retention   int
	signer      Signer
	observers   []Observer
	nextSeq     uint64
	lastHashHex string
}

// AuditOption mutates log configuration.
type AuditOption func(*AuditLog)

// WithRetention limits the maximum number of retained entries.
func WithRetention(limit int) AuditOption {
	return func(log *AuditLog) {
		if limit > 0 {
			log.retention = limit
		}
	}
}

// WithSigner attaches a signing implementation to the audit log.
func WithSigner(s Signer) AuditOption {
	return func(log *AuditLog) {
		log.signer = s
	}
}

// WithObserver registers an observer callback for new entries.
func WithObserver(obs Observer) AuditOption {
	return func(log *AuditLog) {
		if obs != nil {
			log.observers = append(log.observers, obs)
		}
	}
}

// NewAuditLog creates an AuditLog with optional retention and signing.
func NewAuditLog(opts ...AuditOption) *AuditLog {
	log := &AuditLog{retention: 10_000, nextSeq: 1}
	for _, opt := range opts {
		opt(log)
	}
	return log
}

// Append stores an event in the audit log and returns the resulting record.
func (a *AuditLog) Append(event AuditEvent) (AuditRecord, error) {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	if event.ID == "" {
		generated, err := randomID()
		if err != nil {
			return AuditRecord{}, err
		}
		event.ID = generated
	}
	if event.Actor == "" || event.Action == "" {
		return AuditRecord{}, fmt.Errorf("actor and action are required")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	record := AuditRecord{
		Sequence:   a.nextSeq,
		PrevHash:   a.lastHashHex,
		AuditEvent: event,
	}
	payload := digestPayload(record.Sequence, record.PrevHash, event)
	sum := sha256.Sum256(payload)
	record.Hash = hex.EncodeToString(sum[:])
	if a.signer != nil {
		sig, err := a.signer.Sign(sum[:])
		if err != nil {
			return AuditRecord{}, fmt.Errorf("sign audit entry: %w", err)
		}
		record.Signature = sig
	}

	a.entries = append(a.entries, record)
	a.lastHashHex = record.Hash
	a.nextSeq++
	a.enforceRetention()
	a.notify(record)
	return record, nil
}

// Entries returns a copy of log entries in chronological order.
func (a *AuditLog) Entries() []AuditRecord {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]AuditRecord, len(a.entries))
	copy(out, a.entries)
	return out
}

// VerifyChain ensures hash continuity and valid signatures for all entries.
func (a *AuditLog) VerifyChain() error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	prevHash := ""
	for _, record := range a.entries {
		if record.PrevHash != prevHash {
			return fmt.Errorf("hash continuity broken at sequence %d", record.Sequence)
		}
		payload := digestPayload(record.Sequence, record.PrevHash, record.AuditEvent)
		expected := sha256.Sum256(payload)
		if hex.EncodeToString(expected[:]) != record.Hash {
			return fmt.Errorf("hash mismatch at sequence %d", record.Sequence)
		}
		if a.signer != nil && record.Signature != "" {
			if !a.signer.Verify(expected[:], record.Signature) {
				return fmt.Errorf("signature mismatch at sequence %d", record.Sequence)
			}
		}
		prevHash = record.Hash
	}
	return nil
}

// Latest returns the newest record if one exists.
func (a *AuditLog) Latest() (AuditRecord, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.entries) == 0 {
		return AuditRecord{}, false
	}
	return a.entries[len(a.entries)-1], true
}

func (a *AuditLog) enforceRetention() {
	if a.retention <= 0 || len(a.entries) <= a.retention {
		return
	}
	excess := len(a.entries) - a.retention
	a.entries = append([]AuditRecord(nil), a.entries[excess:]...)
}

func (a *AuditLog) notify(record AuditRecord) {
	for _, obs := range a.observers {
		obs(record)
	}
}

func digestPayload(seq uint64, prevHash string, event AuditEvent) []byte {
	metadata := ""
	if len(event.Metadata) > 0 {
		keys := make([]string, 0, len(event.Metadata))
		for k := range event.Metadata {
			keys = append(keys, k)
		}
		// Deterministic ordering ensures consistent digests.
		sort.Strings(keys)
		pairs := make([]string, 0, len(keys))
		for _, k := range keys {
			pairs = append(pairs, k+"="+event.Metadata[k])
		}
		metadata = strings.Join(pairs, ";")
	}
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%d|%s|%s|%s|%s|%s|%d|%s|%s|%s|", seq, prevHash, event.ID, event.Actor, event.Action, event.Scope, event.GasBudget, event.NodeID, event.Reason, metadata))
	builder.WriteString(event.Timestamp.UTC().Format(time.RFC3339Nano))
	return []byte(builder.String())
}

func randomID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate audit id: %w", err)
	}
	return hex.EncodeToString(buf), nil
}
