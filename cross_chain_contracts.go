package synnergy

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// ContractLinkStatus represents the lifecycle state for a cross-chain contract mapping.
type ContractLinkStatus string

const (
	// ContractLinkStatusPending indicates the mapping is awaiting the required approvals.
	ContractLinkStatusPending ContractLinkStatus = "PENDING"
	// ContractLinkStatusActive means the mapping can route traffic across the function web.
	ContractLinkStatusActive ContractLinkStatus = "ACTIVE"
	// ContractLinkStatusSuspended means the mapping has been temporarily disabled.
	ContractLinkStatusSuspended ContractLinkStatus = "SUSPENDED"
	// ContractLinkStatusRetired means the mapping has been permanently decommissioned.
	ContractLinkStatusRetired ContractLinkStatus = "RETIRED"
	// ContractLinkStatusFailed captures a mapping that faulted and needs remediation.
	ContractLinkStatusFailed ContractLinkStatus = "FAILED"
)

// ContractLinkEventType describes the reason an event was emitted.
type ContractLinkEventType string

const (
	ContractLinkEventRegistered       ContractLinkEventType = "REGISTERED"
	ContractLinkEventApprovalRecorded ContractLinkEventType = "APPROVAL_RECORDED"
	ContractLinkEventActivated        ContractLinkEventType = "ACTIVATED"
	ContractLinkEventSuspended        ContractLinkEventType = "SUSPENDED"
	ContractLinkEventResumed          ContractLinkEventType = "RESUMED"
	ContractLinkEventUpdated          ContractLinkEventType = "UPDATED"
	ContractLinkEventRetired          ContractLinkEventType = "RETIRED"
	ContractLinkEventFailed           ContractLinkEventType = "FAILED"
)

// AccessPolicy defines the approval and privacy controls applied to a mapping.
type AccessPolicy struct {
	AllowedApprovers  []string
	RequiredApprovals int
	PrivacyLevel      string
	EncryptionScheme  string
}

// ContractLinkSpec captures the configuration for a cross-chain contract mapping.
type ContractLinkSpec struct {
	LocalChain     string
	LocalAddress   string
	RemoteChain    string
	RemoteAddress  string
	ConnectionID   string
	Capabilities   []string
	GasProfile     string
	Metadata       map[string]string
	AccessPolicy   AccessPolicy
	AuditTrailHint string
}

// ContractLinkFailure records failure context for a contract mapping.
type ContractLinkFailure struct {
	Code       string
	Detail     string
	Occurred   time.Time
	Resolved   bool
	ResolvedAt time.Time
}

// ContractLink models the lifecycle of a cross-chain contract mapping.
type ContractLink struct {
	ID                string
	Spec              ContractLinkSpec
	Status            ContractLinkStatus
	Version           uint64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ActivatedAt       time.Time
	SuspendedAt       time.Time
	RetiredAt         time.Time
	SuspensionReason  string
	RetirementReason  string
	ApprovalThreshold int
	Approvals         map[string]time.Time
	Failures          []ContractLinkFailure
}

// ContractLinkEvent notifies subscribers of lifecycle changes.
type ContractLinkEvent struct {
	Type     ContractLinkEventType
	LinkID   string
	Link     *ContractLink
	Approver string
	Reason   string
	Failure  *ContractLinkFailure
}

// ContractLinkFilter allows callers to constrain list results.
type ContractLinkFilter struct {
	Statuses     []ContractLinkStatus
	LocalChain   string
	RemoteChain  string
	ConnectionID string
}

// ConnectionSnapshotSource provides read access to connection snapshots.
type ConnectionSnapshotSource interface {
	GetConnection(id string) (*ChainConnection, bool)
}

// ContractLinkIDGenerator produces deterministic identifiers for mappings.
type ContractLinkIDGenerator interface {
	NewContractLinkID(spec ContractLinkSpec) (string, error)
}

// defaultContractLinkIDGenerator hashes the specification to produce a stable ID.
type defaultContractLinkIDGenerator struct{}

// NewContractLinkID implements ContractLinkIDGenerator.
func (defaultContractLinkIDGenerator) NewContractLinkID(spec ContractLinkSpec) (string, error) {
	segments := []string{
		strings.ToLower(strings.TrimSpace(spec.LocalChain)),
		strings.ToLower(strings.TrimSpace(spec.LocalAddress)),
		strings.ToLower(strings.TrimSpace(spec.RemoteChain)),
		strings.ToLower(strings.TrimSpace(spec.RemoteAddress)),
		strings.ToLower(strings.TrimSpace(spec.ConnectionID)),
	}
	for _, segment := range segments {
		if segment == "" {
			return "", fmt.Errorf("%w: specification missing required fields", ErrContractLinkInvalidSpec)
		}
	}
	sum := sha256.Sum256([]byte(strings.Join(segments, "|")))
	return hex.EncodeToString(sum[:]), nil
}

// ContractLinkManager orchestrates contract link lifecycles and telemetry.
type ContractLinkManager struct {
	mu       sync.RWMutex
	links    map[string]*ContractLink
	byLocal  map[string]string
	watchers map[int]chan ContractLinkEvent
	nextID   int
	resolver ConnectionSnapshotSource
	idGen    ContractLinkIDGenerator
	clock    func() time.Time
}

// ContractLinkManagerOption customises manager construction.
type ContractLinkManagerOption func(*ContractLinkManager)

// WithContractLinkResolver configures the manager to verify connections.
func WithContractLinkResolver(resolver ConnectionSnapshotSource) ContractLinkManagerOption {
	return func(m *ContractLinkManager) {
		m.resolver = resolver
	}
}

// WithContractLinkIDGenerator overrides the identifier generator.
func WithContractLinkIDGenerator(gen ContractLinkIDGenerator) ContractLinkManagerOption {
	return func(m *ContractLinkManager) {
		if gen != nil {
			m.idGen = gen
		}
	}
}

// WithContractLinkClock overrides the clock used for timestamps (primarily for tests).
func WithContractLinkClock(clock func() time.Time) ContractLinkManagerOption {
	return func(m *ContractLinkManager) {
		if clock != nil {
			m.clock = clock
		}
	}
}

// NewContractLinkManager builds a hardened manager suitable for enterprise workflows.
func NewContractLinkManager(opts ...ContractLinkManagerOption) *ContractLinkManager {
	m := &ContractLinkManager{
		links:    make(map[string]*ContractLink),
		byLocal:  make(map[string]string),
		watchers: make(map[int]chan ContractLinkEvent),
		idGen:    defaultContractLinkIDGenerator{},
		clock: func() time.Time {
			return time.Now().UTC()
		},
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

var (
	// ErrContractLinkExists indicates the mapping already exists.
	ErrContractLinkExists = errors.New("contract link already exists")
	// ErrContractLinkNotFound indicates the requested mapping does not exist.
	ErrContractLinkNotFound = errors.New("contract link not found")
	// ErrContractLinkConnectionUnknown indicates the referenced connection is missing.
	ErrContractLinkConnectionUnknown = errors.New("contract link connection unknown")
	// ErrContractLinkConnectionInactive indicates the referenced connection is not active.
	ErrContractLinkConnectionInactive = errors.New("contract link connection inactive")
	// ErrContractLinkInvalidSpec indicates the provided specification is invalid.
	ErrContractLinkInvalidSpec = errors.New("contract link specification invalid")
	// ErrContractLinkApproverNotAllowed indicates the approver is not authorised for the mapping.
	ErrContractLinkApproverNotAllowed = errors.New("approver not permitted for this contract link")
	// ErrContractLinkApprovalDuplicate indicates the approver has already signed off on the mapping.
	ErrContractLinkApprovalDuplicate = errors.New("approval already recorded for contract link")
	// ErrContractLinkInvalidState indicates the mapping is not in a state that supports the requested transition.
	ErrContractLinkInvalidState = errors.New("contract link is in an invalid state for this operation")
)

// ContractLinkUpdate captures mutable attributes for a mapping.
type ContractLinkUpdate struct {
	Capabilities []string
	GasProfile   string
	Metadata     map[string]string
}

// Register adds a new contract mapping into the catalogue.
func (m *ContractLinkManager) Register(spec ContractLinkSpec) (*ContractLink, error) {
	if err := validateContractLinkSpec(spec); err != nil {
		return nil, err
	}
	if err := m.ensureConnectionActive(spec.ConnectionID); err != nil {
		return nil, err
	}
	canonicalSpec := cloneContractLinkSpec(spec)
	id, err := m.idGen.NewContractLinkID(canonicalSpec)
	if err != nil {
		return nil, err
	}
	now := m.now()
	approvals := make(map[string]time.Time)
	status := ContractLinkStatusPending
	if canonicalSpec.AccessPolicy.RequiredApprovals <= 0 {
		status = ContractLinkStatusActive
	}
	link := &ContractLink{
		ID:                id,
		Spec:              canonicalSpec,
		Status:            status,
		Version:           1,
		CreatedAt:         now,
		UpdatedAt:         now,
		ApprovalThreshold: canonicalSpec.AccessPolicy.RequiredApprovals,
		Approvals:         approvals,
	}
	if status == ContractLinkStatusActive {
		link.ActivatedAt = now
	}
	localKey := normalizeKey(canonicalSpec.LocalChain, canonicalSpec.LocalAddress)

	m.mu.Lock()
	if _, exists := m.links[id]; exists {
		m.mu.Unlock()
		return nil, ErrContractLinkExists
	}
	if existing, ok := m.byLocal[localKey]; ok {
		m.mu.Unlock()
		return nil, fmt.Errorf("%w: local contract already mapped to %s", ErrContractLinkExists, existing)
	}
	m.links[id] = link
	m.byLocal[localKey] = id
	snapshot := cloneContractLink(link)
	activated := status == ContractLinkStatusActive
	m.mu.Unlock()

	m.broadcast(ContractLinkEvent{Type: ContractLinkEventRegistered, LinkID: id, Link: snapshot})
	if activated {
		m.broadcast(ContractLinkEvent{Type: ContractLinkEventActivated, LinkID: id, Link: snapshot})
	}
	return snapshot, nil
}

// RecordApproval records an approval for a mapping and activates it when the threshold is reached.
func (m *ContractLinkManager) RecordApproval(linkID, approver string) (*ContractLink, error) {
	normalizedApprover := strings.ToLower(strings.TrimSpace(approver))
	if normalizedApprover == "" {
		return nil, fmt.Errorf("%w: approver required", ErrContractLinkApproverNotAllowed)
	}
	m.mu.Lock()
	link, ok := m.links[linkID]
	if !ok {
		m.mu.Unlock()
		return nil, ErrContractLinkNotFound
	}
	if link.Status != ContractLinkStatusPending {
		m.mu.Unlock()
		return nil, ErrContractLinkInvalidState
	}
	if !approverAllowed(link.Spec.AccessPolicy, approver) {
		m.mu.Unlock()
		return nil, ErrContractLinkApproverNotAllowed
	}
	if _, exists := link.Approvals[normalizedApprover]; exists {
		m.mu.Unlock()
		return nil, ErrContractLinkApprovalDuplicate
	}
	now := m.now()
	link.Approvals[normalizedApprover] = now
	link.UpdatedAt = now
	link.Version++
	activated := false
	if len(link.Approvals) >= link.ApprovalThreshold {
		link.Status = ContractLinkStatusActive
		link.ActivatedAt = now
		activated = true
	}
	snapshot := cloneContractLink(link)
	m.mu.Unlock()

	m.broadcast(ContractLinkEvent{Type: ContractLinkEventApprovalRecorded, LinkID: linkID, Link: snapshot, Approver: approver})
	if activated {
		m.broadcast(ContractLinkEvent{Type: ContractLinkEventActivated, LinkID: linkID, Link: snapshot})
	}
	return snapshot, nil
}

// Update mutates metadata, capabilities, or gas profile for a mapping.
func (m *ContractLinkManager) Update(linkID string, update ContractLinkUpdate) (*ContractLink, error) {
	m.mu.Lock()
	link, ok := m.links[linkID]
	if !ok {
		m.mu.Unlock()
		return nil, ErrContractLinkNotFound
	}
	if link.Status == ContractLinkStatusRetired {
		m.mu.Unlock()
		return nil, ErrContractLinkInvalidState
	}
	changed := false
	if update.GasProfile != "" && update.GasProfile != link.Spec.GasProfile {
		link.Spec.GasProfile = strings.TrimSpace(update.GasProfile)
		changed = true
	}
	if update.Metadata != nil {
		link.Spec.Metadata = cloneStringMap(update.Metadata)
		changed = true
	}
	if update.Capabilities != nil {
		link.Spec.Capabilities = cloneStringSlice(update.Capabilities)
		changed = true
	}
	if changed {
		link.Version++
		link.UpdatedAt = m.now()
	}
	snapshot := cloneContractLink(link)
	m.mu.Unlock()
	if changed {
		m.broadcast(ContractLinkEvent{Type: ContractLinkEventUpdated, LinkID: linkID, Link: snapshot})
	}
	return snapshot, nil
}

// Suspend disables a mapping while retaining its configuration.
func (m *ContractLinkManager) Suspend(linkID, reason string) (*ContractLink, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, fmt.Errorf("%w: suspension reason required", ErrContractLinkInvalidState)
	}
	m.mu.Lock()
	link, ok := m.links[linkID]
	if !ok {
		m.mu.Unlock()
		return nil, ErrContractLinkNotFound
	}
	if link.Status != ContractLinkStatusActive {
		m.mu.Unlock()
		return nil, ErrContractLinkInvalidState
	}
	now := m.now()
	link.Status = ContractLinkStatusSuspended
	link.SuspensionReason = reason
	link.SuspendedAt = now
	link.UpdatedAt = now
	link.Version++
	snapshot := cloneContractLink(link)
	m.mu.Unlock()
	m.broadcast(ContractLinkEvent{Type: ContractLinkEventSuspended, LinkID: linkID, Link: snapshot, Reason: reason})
	return snapshot, nil
}

// Resume re-enables a suspended or failed mapping after remediation.
func (m *ContractLinkManager) Resume(linkID string) (*ContractLink, error) {
	m.mu.Lock()
	link, ok := m.links[linkID]
	if !ok {
		m.mu.Unlock()
		return nil, ErrContractLinkNotFound
	}
	if link.Status != ContractLinkStatusSuspended && link.Status != ContractLinkStatusFailed {
		m.mu.Unlock()
		return nil, ErrContractLinkInvalidState
	}
	now := m.now()
	link.Status = ContractLinkStatusActive
	link.SuspensionReason = ""
	link.SuspendedAt = time.Time{}
	link.UpdatedAt = now
	link.Version++
	for i := len(link.Failures) - 1; i >= 0; i-- {
		if !link.Failures[i].Resolved {
			link.Failures[i].Resolved = true
			link.Failures[i].ResolvedAt = now
			break
		}
	}
	snapshot := cloneContractLink(link)
	m.mu.Unlock()
	m.broadcast(ContractLinkEvent{Type: ContractLinkEventResumed, LinkID: linkID, Link: snapshot})
	return snapshot, nil
}

// Retire decommissions a mapping permanently.
func (m *ContractLinkManager) Retire(linkID, reason string) (*ContractLink, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, fmt.Errorf("%w: retirement reason required", ErrContractLinkInvalidState)
	}
	m.mu.Lock()
	link, ok := m.links[linkID]
	if !ok {
		m.mu.Unlock()
		return nil, ErrContractLinkNotFound
	}
	if link.Status == ContractLinkStatusRetired {
		m.mu.Unlock()
		return nil, ErrContractLinkInvalidState
	}
	now := m.now()
	link.Status = ContractLinkStatusRetired
	link.RetirementReason = reason
	link.RetiredAt = now
	link.UpdatedAt = now
	link.Version++
	localKey := normalizeKey(link.Spec.LocalChain, link.Spec.LocalAddress)
	delete(m.byLocal, localKey)
	snapshot := cloneContractLink(link)
	m.mu.Unlock()
	m.broadcast(ContractLinkEvent{Type: ContractLinkEventRetired, LinkID: linkID, Link: snapshot, Reason: reason})
	return snapshot, nil
}

// ReportFailure records a fault and transitions the mapping to FAILED.
func (m *ContractLinkManager) ReportFailure(linkID, code, detail string) (*ContractLink, error) {
	code = strings.TrimSpace(code)
	detail = strings.TrimSpace(detail)
	if code == "" || detail == "" {
		return nil, fmt.Errorf("%w: failure code and detail required", ErrContractLinkInvalidState)
	}
	m.mu.Lock()
	link, ok := m.links[linkID]
	if !ok {
		m.mu.Unlock()
		return nil, ErrContractLinkNotFound
	}
	if link.Status == ContractLinkStatusRetired {
		m.mu.Unlock()
		return nil, ErrContractLinkInvalidState
	}
	now := m.now()
	failure := ContractLinkFailure{Code: code, Detail: detail, Occurred: now}
	link.Failures = append(link.Failures, failure)
	link.Status = ContractLinkStatusFailed
	link.SuspensionReason = detail
	link.SuspendedAt = now
	link.UpdatedAt = now
	link.Version++
	snapshot := cloneContractLink(link)
	m.mu.Unlock()
	m.broadcast(ContractLinkEvent{Type: ContractLinkEventFailed, LinkID: linkID, Link: snapshot, Reason: detail, Failure: &failure})
	return snapshot, nil
}

// Get returns a copy of the mapping by ID.
func (m *ContractLinkManager) Get(linkID string) (*ContractLink, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	link, ok := m.links[linkID]
	if !ok {
		return nil, false
	}
	return cloneContractLink(link), true
}

// List enumerates known mappings according to the provided filter.
func (m *ContractLinkManager) List(filter ContractLinkFilter) []*ContractLink {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var statuses map[ContractLinkStatus]struct{}
	if len(filter.Statuses) > 0 {
		statuses = make(map[ContractLinkStatus]struct{}, len(filter.Statuses))
		for _, st := range filter.Statuses {
			statuses[st] = struct{}{}
		}
	}
	var out []*ContractLink
	for _, link := range m.links {
		if statuses != nil {
			if _, ok := statuses[link.Status]; !ok {
				continue
			}
		}
		if filter.LocalChain != "" && !strings.EqualFold(filter.LocalChain, link.Spec.LocalChain) {
			continue
		}
		if filter.RemoteChain != "" && !strings.EqualFold(filter.RemoteChain, link.Spec.RemoteChain) {
			continue
		}
		if filter.ConnectionID != "" && filter.ConnectionID != link.Spec.ConnectionID {
			continue
		}
		out = append(out, cloneContractLink(link))
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out
}

// Subscribe registers a watcher for lifecycle events.
func (m *ContractLinkManager) Subscribe(buffer int) (<-chan ContractLinkEvent, func()) {
	if buffer <= 0 {
		buffer = 1
	}
	m.mu.Lock()
	id := m.nextID
	m.nextID++
	ch := make(chan ContractLinkEvent, buffer)
	m.watchers[id] = ch
	m.mu.Unlock()
	cancel := func() {
		m.mu.Lock()
		if watcher, ok := m.watchers[id]; ok {
			delete(m.watchers, id)
			close(watcher)
		}
		m.mu.Unlock()
	}
	return ch, cancel
}

func (m *ContractLinkManager) broadcast(event ContractLinkEvent) {
	m.mu.RLock()
	watchers := make([]chan ContractLinkEvent, 0, len(m.watchers))
	for _, ch := range m.watchers {
		watchers = append(watchers, ch)
	}
	m.mu.RUnlock()
	for _, ch := range watchers {
		select {
		case ch <- event:
		default:
		}
	}
}

func (m *ContractLinkManager) ensureConnectionActive(connectionID string) error {
	if strings.TrimSpace(connectionID) == "" {
		return fmt.Errorf("%w: connection id required", ErrContractLinkInvalidSpec)
	}
	if m.resolver == nil {
		return nil
	}
	conn, ok := m.resolver.GetConnection(connectionID)
	if !ok {
		return ErrContractLinkConnectionUnknown
	}
	if conn.Status != ConnectionStatusActive {
		return ErrContractLinkConnectionInactive
	}
	return nil
}

func (m *ContractLinkManager) now() time.Time {
	if m.clock != nil {
		return m.clock()
	}
	return time.Now().UTC()
}

func normalizeKey(chain, address string) string {
	return strings.ToLower(strings.TrimSpace(chain)) + "|" + strings.ToLower(strings.TrimSpace(address))
}

func cloneContractLink(link *ContractLink) *ContractLink {
	if link == nil {
		return nil
	}
	approvals := make(map[string]time.Time, len(link.Approvals))
	for k, v := range link.Approvals {
		approvals[k] = v
	}
	failures := make([]ContractLinkFailure, len(link.Failures))
	copy(failures, link.Failures)
	return &ContractLink{
		ID:                link.ID,
		Spec:              cloneContractLinkSpec(link.Spec),
		Status:            link.Status,
		Version:           link.Version,
		CreatedAt:         link.CreatedAt,
		UpdatedAt:         link.UpdatedAt,
		ActivatedAt:       link.ActivatedAt,
		SuspendedAt:       link.SuspendedAt,
		RetiredAt:         link.RetiredAt,
		SuspensionReason:  link.SuspensionReason,
		RetirementReason:  link.RetirementReason,
		ApprovalThreshold: link.ApprovalThreshold,
		Approvals:         approvals,
		Failures:          failures,
	}
}

func cloneContractLinkSpec(spec ContractLinkSpec) ContractLinkSpec {
	return ContractLinkSpec{
		LocalChain:    strings.TrimSpace(spec.LocalChain),
		LocalAddress:  strings.TrimSpace(spec.LocalAddress),
		RemoteChain:   strings.TrimSpace(spec.RemoteChain),
		RemoteAddress: strings.TrimSpace(spec.RemoteAddress),
		ConnectionID:  strings.TrimSpace(spec.ConnectionID),
		Capabilities:  cloneStringSlice(spec.Capabilities),
		GasProfile:    strings.TrimSpace(spec.GasProfile),
		Metadata:      cloneStringMap(spec.Metadata),
		AccessPolicy: AccessPolicy{
			AllowedApprovers:  cloneAndSortApprovers(spec.AccessPolicy.AllowedApprovers),
			RequiredApprovals: spec.AccessPolicy.RequiredApprovals,
			PrivacyLevel:      strings.TrimSpace(spec.AccessPolicy.PrivacyLevel),
			EncryptionScheme:  strings.TrimSpace(spec.AccessPolicy.EncryptionScheme),
		},
		AuditTrailHint: strings.TrimSpace(spec.AuditTrailHint),
	}
}

func cloneStringSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	out := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func cloneAndSortApprovers(values []string) []string {
	out := cloneStringSlice(values)
	if len(out) == 0 {
		return nil
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i]) < strings.ToLower(out[j])
	})
	return out
}

func cloneStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		key := strings.TrimSpace(k)
		if key == "" {
			continue
		}
		out[key] = strings.TrimSpace(v)
	}
	return out
}

func approverAllowed(policy AccessPolicy, approver string) bool {
	if len(policy.AllowedApprovers) == 0 {
		return true
	}
	candidate := strings.TrimSpace(approver)
	for _, allowed := range policy.AllowedApprovers {
		if strings.EqualFold(allowed, candidate) {
			return true
		}
	}
	return false
}

func validateContractLinkSpec(spec ContractLinkSpec) error {
	if strings.TrimSpace(spec.LocalChain) == "" {
		return fmt.Errorf("%w: local chain required", ErrContractLinkInvalidSpec)
	}
	if strings.TrimSpace(spec.LocalAddress) == "" {
		return fmt.Errorf("%w: local address required", ErrContractLinkInvalidSpec)
	}
	if strings.TrimSpace(spec.RemoteChain) == "" {
		return fmt.Errorf("%w: remote chain required", ErrContractLinkInvalidSpec)
	}
	if strings.TrimSpace(spec.RemoteAddress) == "" {
		return fmt.Errorf("%w: remote address required", ErrContractLinkInvalidSpec)
	}
	if strings.TrimSpace(spec.ConnectionID) == "" {
		return fmt.Errorf("%w: connection id required", ErrContractLinkInvalidSpec)
	}
	if spec.AccessPolicy.RequiredApprovals < 0 {
		return fmt.Errorf("%w: required approvals cannot be negative", ErrContractLinkInvalidSpec)
	}
	if spec.AccessPolicy.RequiredApprovals > 0 {
		if len(spec.AccessPolicy.AllowedApprovers) == 0 {
			return fmt.Errorf("%w: allowed approvers required when approvals are mandated", ErrContractLinkInvalidSpec)
		}
		unique := make(map[string]struct{}, len(spec.AccessPolicy.AllowedApprovers))
		for _, approver := range spec.AccessPolicy.AllowedApprovers {
			key := strings.ToLower(strings.TrimSpace(approver))
			if key == "" {
				continue
			}
			unique[key] = struct{}{}
		}
		if spec.AccessPolicy.RequiredApprovals > len(unique) {
			return fmt.Errorf("%w: approval threshold exceeds distinct approvers", ErrContractLinkInvalidSpec)
		}
	}
	return nil
}
