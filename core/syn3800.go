package core

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// GrantStatus captures the lifecycle stage for a SYN3800 grant token.
type GrantStatus string

const (
	// GrantStatusPending indicates the grant has been created but no funds were
	// released yet.
	GrantStatusPending GrantStatus = "PENDING"
	// GrantStatusActive marks grants that have partially released funds.
	GrantStatusActive GrantStatus = "ACTIVE"
	// GrantStatusCompleted marks grants whose allocation has been
	// completely disbursed.
	GrantStatusCompleted GrantStatus = "COMPLETED"
	// GrantStatusRevoked marks grants that were revoked for risk or
	// compliance reasons.
	GrantStatusRevoked GrantStatus = "REVOKED"
)

// GrantEventType enumerates audit events recorded for each grant.
type GrantEventType string

const (
	GrantEventCreated    GrantEventType = "CREATED"
	GrantEventAuthorized GrantEventType = "AUTHORIZED"
	GrantEventRevoked    GrantEventType = "REVOKED"
	GrantEventDisbursed  GrantEventType = "DISBURSED"
	GrantEventNote       GrantEventType = "NOTE"
)

var (
	// ErrGrantNotFound is returned when the requested grant ID is missing.
	ErrGrantNotFound = errors.New("grant not found")
	// ErrGrantInvalidAmount is returned when a disbursement amount is zero or
	// exceeds the remaining allocation.
	ErrGrantInvalidAmount = errors.New("invalid disbursement amount")
	// ErrGrantUnauthorized indicates the acting wallet is not authorized to
	// operate on the grant.
	ErrGrantUnauthorized = errors.New("wallet not authorised")
)

// GrantEvent represents a single lifecycle event for a grant record.
type GrantEvent struct {
	Timestamp time.Time      `json:"timestamp"`
	Type      GrantEventType `json:"type"`
	Signer    string         `json:"signer,omitempty"`
	Amount    uint64         `json:"amount,omitempty"`
	Note      string         `json:"note,omitempty"`
}

// GrantRecord captures metadata, authorizers and audit trail for a SYN3800 grant token.
type GrantRecord struct {
	ID          uint64               `json:"id"`
	Beneficiary string               `json:"beneficiary"`
	Name        string               `json:"name"`
	Amount      uint64               `json:"amount"`
	Released    uint64               `json:"released"`
	Notes       []string             `json:"notes,omitempty"`
	Status      GrantStatus          `json:"status"`
	Authorizers map[string]time.Time `json:"authorizers,omitempty"`
	Events      []GrantEvent         `json:"events,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	RevokedBy   string               `json:"revoked_by,omitempty"`
	RevokeNote  string               `json:"revoke_note,omitempty"`
}

func (g *GrantRecord) clone() *GrantRecord {
	if g == nil {
		return nil
	}
	cp := *g
	if len(g.Notes) > 0 {
		cp.Notes = append([]string(nil), g.Notes...)
	}
	if len(g.Events) > 0 {
		cp.Events = append([]GrantEvent(nil), g.Events...)
	}
	if len(g.Authorizers) > 0 {
		cp.Authorizers = make(map[string]time.Time, len(g.Authorizers))
		for k, v := range g.Authorizers {
			cp.Authorizers[k] = v
		}
	}
	return &cp
}

// GrantRegistrySnapshot represents the persisted state for Stage 73 storage.
type GrantRegistrySnapshot struct {
	NextID uint64         `json:"next_id"`
	Grants []*GrantRecord `json:"grants"`
}

// GrantRegistry manages grant records.
type GrantRegistry struct {
	mu     sync.RWMutex
	grants map[uint64]*GrantRecord
	nextID uint64
}

// NewGrantRegistry creates a new registry.
func NewGrantRegistry() *GrantRegistry {
	return &GrantRegistry{grants: make(map[uint64]*GrantRecord)}
}

// CreateGrant registers a new grant, optional authorizers and returns its ID.
func (r *GrantRegistry) CreateGrant(beneficiary, name string, amount uint64, authorizers ...string) (uint64, error) {
	if beneficiary == "" {
		return 0, errors.New("beneficiary required")
	}
	if name == "" {
		return 0, errors.New("name required")
	}
	if amount == 0 {
		return 0, errors.New("invalid amount")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	id := r.nextID
	now := time.Now().UTC()
	auths := make(map[string]time.Time)
	for _, addr := range authorizers {
		if addr == "" {
			continue
		}
		auths[addr] = now
	}
	record := &GrantRecord{
		ID:          id,
		Beneficiary: beneficiary,
		Name:        name,
		Amount:      amount,
		Status:      GrantStatusPending,
		Authorizers: auths,
		CreatedAt:   now,
		UpdatedAt:   now,
		Events: []GrantEvent{{
			Timestamp: now,
			Type:      GrantEventCreated,
			Note:      fmt.Sprintf("created for %s", beneficiary),
		}},
	}
	if len(auths) > 0 {
		record.Events = append(record.Events, GrantEvent{
			Timestamp: now,
			Type:      GrantEventAuthorized,
			Note:      "initial authorizers",
		})
	}
	r.grants[id] = record
	return id, nil
}

// Authorize registers a new wallet address authorised to operate on the grant.
func (r *GrantRegistry) Authorize(id uint64, addr string) (*GrantEvent, error) {
	if addr == "" {
		return nil, errors.New("address required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.grants[id]
	if !ok {
		return nil, ErrGrantNotFound
	}
	if g.Authorizers == nil {
		g.Authorizers = make(map[string]time.Time)
	}
	now := time.Now().UTC()
	if _, exists := g.Authorizers[addr]; exists {
		// Even if already authorized, record a note so audit trail remains complete.
		evt := GrantEvent{Timestamp: now, Type: GrantEventNote, Signer: addr, Note: "re-authorized"}
		g.Events = append(g.Events, evt)
		g.UpdatedAt = now
		return &evt, nil
	}
	g.Authorizers[addr] = now
	evt := GrantEvent{Timestamp: now, Type: GrantEventAuthorized, Signer: addr}
	g.Events = append(g.Events, evt)
	g.UpdatedAt = now
	return &evt, nil
}

// Revoke removes an address from the authorised set and records the reason.
func (r *GrantRegistry) Revoke(id uint64, addr, note string) (*GrantEvent, error) {
	if addr == "" {
		return nil, errors.New("address required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.grants[id]
	if !ok {
		return nil, ErrGrantNotFound
	}
	if g.Authorizers == nil {
		return nil, ErrGrantUnauthorized
	}
	if _, exists := g.Authorizers[addr]; !exists {
		return nil, ErrGrantUnauthorized
	}
	delete(g.Authorizers, addr)
	now := time.Now().UTC()
	evt := GrantEvent{Timestamp: now, Type: GrantEventRevoked, Signer: addr, Note: note}
	g.Events = append(g.Events, evt)
	g.Status = GrantStatusRevoked
	g.RevokedBy = addr
	g.RevokeNote = note
	g.UpdatedAt = now
	return &evt, nil
}

// Disburse releases a portion of the grant and records an audit event.
func (r *GrantRegistry) Disburse(id uint64, amount uint64, note, signer string) (*GrantEvent, error) {
	if amount == 0 {
		return nil, ErrGrantInvalidAmount
	}
	if signer == "" {
		return nil, ErrGrantUnauthorized
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.grants[id]
	if !ok {
		return nil, ErrGrantNotFound
	}
	if g.Status == GrantStatusRevoked {
		return nil, errors.New("grant revoked")
	}
	if g.Authorizers != nil {
		if _, exists := g.Authorizers[signer]; !exists {
			return nil, ErrGrantUnauthorized
		}
	} else {
		return nil, ErrGrantUnauthorized
	}
	if g.Released+amount > g.Amount {
		return nil, ErrGrantInvalidAmount
	}
	g.Released += amount
	if note != "" {
		g.Notes = append(g.Notes, note)
	}
	now := time.Now().UTC()
	evt := GrantEvent{Timestamp: now, Type: GrantEventDisbursed, Signer: signer, Amount: amount, Note: note}
	g.Events = append(g.Events, evt)
	if g.Released == g.Amount {
		g.Status = GrantStatusCompleted
	} else {
		g.Status = GrantStatusActive
	}
	g.UpdatedAt = now
	return &evt, nil
}

// GetGrant returns a clone of the grant record by ID.
func (r *GrantRegistry) GetGrant(id uint64) (*GrantRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.grants[id]
	if !ok {
		return nil, false
	}
	return g.clone(), true
}

// ListGrants returns sorted copies of all grants for deterministic telemetry and API output.
func (r *GrantRegistry) ListGrants() []*GrantRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*GrantRecord, 0, len(r.grants))
	for _, g := range r.grants {
		res = append(res, g.clone())
	}
	sort.Slice(res, func(i, j int) bool { return res[i].ID < res[j].ID })
	return res
}

// Audit returns a copy of the audit trail for the given grant.
func (r *GrantRegistry) Audit(id uint64) ([]GrantEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.grants[id]
	if !ok {
		return nil, ErrGrantNotFound
	}
	events := append([]GrantEvent(nil), g.Events...)
	sort.Slice(events, func(i, j int) bool { return events[i].Timestamp.Before(events[j].Timestamp) })
	return events, nil
}

// GrantStatusSummary aggregates lifecycle counts for telemetry and reporting.
type GrantStatusSummary struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Active    int `json:"active"`
	Completed int `json:"completed"`
	Revoked   int `json:"revoked"`
}

// StatusSummary calculates totals for each grant status.
func (r *GrantRegistry) StatusSummary() GrantStatusSummary {
	r.mu.RLock()
	defer r.mu.RUnlock()
	summary := GrantStatusSummary{}
	for _, g := range r.grants {
		summary.Total++
		switch g.Status {
		case GrantStatusActive:
			summary.Active++
		case GrantStatusCompleted:
			summary.Completed++
		case GrantStatusRevoked:
			summary.Revoked++
		default:
			summary.Pending++
		}
	}
	return summary
}

// Snapshot returns a deep copy snapshot suitable for persistence.
func (r *GrantRegistry) Snapshot() GrantRegistrySnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	snapshot := GrantRegistrySnapshot{NextID: r.nextID}
	snapshot.Grants = make([]*GrantRecord, 0, len(r.grants))
	for _, g := range r.grants {
		snapshot.Grants = append(snapshot.Grants, g.clone())
	}
	sort.Slice(snapshot.Grants, func(i, j int) bool { return snapshot.Grants[i].ID < snapshot.Grants[j].ID })
	return snapshot
}

// Restore replaces the registry state with the provided snapshot.
func (r *GrantRegistry) Restore(snapshot GrantRegistrySnapshot) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.grants = make(map[uint64]*GrantRecord, len(snapshot.Grants))
	for _, g := range snapshot.Grants {
		if g == nil {
			continue
		}
		clone := g.clone()
		r.grants[clone.ID] = clone
	}
	r.nextID = snapshot.NextID
}

// IsAuthorized reports whether the provided address is authorised for the grant.
func (r *GrantRegistry) IsAuthorized(id uint64, addr string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.grants[id]
	if !ok || addr == "" {
		return false
	}
	if g.Authorizers == nil {
		return false
	}
	_, ok = g.Authorizers[addr]
	return ok
}
