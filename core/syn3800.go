package core

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// GrantStatus represents the lifecycle of a grant.
type GrantStatus string

const (
	// GrantStatusPending indicates the grant exists but no disbursement has occurred yet.
	GrantStatusPending GrantStatus = "PENDING"
	// GrantStatusActive indicates partial disbursements are in flight.
	GrantStatusActive GrantStatus = "ACTIVE"
	// GrantStatusCompleted indicates the grant was fully released.
	GrantStatusCompleted GrantStatus = "COMPLETED"
	// GrantStatusRevoked indicates the grant was revoked before completion.
	GrantStatusRevoked GrantStatus = "REVOKED"
)

// GrantEvent captures a lifecycle transition or disbursement detail for audit trails.
type GrantEvent struct {
	Timestamp time.Time   `json:"timestamp"`
	Actor     string      `json:"actor"`
	Action    string      `json:"action"`
	Amount    uint64      `json:"amount,omitempty"`
	Note      string      `json:"note,omitempty"`
	Status    GrantStatus `json:"status"`
}

// GrantRecord captures metadata and audit history for a SYN3800 grant token.
type GrantRecord struct {
	ID          uint64          `json:"id"`
	Beneficiary string          `json:"beneficiary"`
	Name        string          `json:"name"`
	Amount      uint64          `json:"amount"`
	Released    uint64          `json:"released"`
	Notes       []string        `json:"notes,omitempty"`
	Status      GrantStatus     `json:"status"`
	Authorizers map[string]bool `json:"authorizers"`
	Events      []GrantEvent    `json:"events,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// clone returns a deep copy to preserve immutability for callers.
func (g *GrantRecord) clone() *GrantRecord {
	cp := *g
	if len(g.Notes) > 0 {
		cp.Notes = append([]string{}, g.Notes...)
	}
	if len(g.Authorizers) > 0 {
		cp.Authorizers = make(map[string]bool, len(g.Authorizers))
		for k, v := range g.Authorizers {
			cp.Authorizers[k] = v
		}
	}
	if len(g.Events) > 0 {
		cp.Events = append([]GrantEvent{}, g.Events...)
	}
	return &cp
}

// GrantRegistry manages grant records and their audit trails.
type GrantRegistry struct {
	mu     sync.RWMutex
	grants map[uint64]*GrantRecord
	nextID uint64
}

// NewGrantRegistry creates a new registry.
func NewGrantRegistry() *GrantRegistry {
	return &GrantRegistry{grants: make(map[uint64]*GrantRecord)}
}

// NewGrantRegistryFromRecords restores a registry from persisted records.
func NewGrantRegistryFromRecords(records []*GrantRecord) *GrantRegistry {
	reg := NewGrantRegistry()
	var maxID uint64
	for _, rec := range records {
		if rec == nil {
			continue
		}
		cp := rec.clone()
		reg.grants[cp.ID] = cp
		if cp.ID > maxID {
			maxID = cp.ID
		}
	}
	reg.nextID = maxID
	return reg
}

// CreateGrant registers a new grant and returns its ID.
func (r *GrantRegistry) CreateGrant(beneficiary, name string, amount uint64, authorizer string) uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	id := r.nextID
	now := time.Now().UTC()
	rec := &GrantRecord{
		ID:          id,
		Beneficiary: beneficiary,
		Name:        name,
		Amount:      amount,
		Status:      GrantStatusPending,
		Authorizers: make(map[string]bool),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if authorizer != "" {
		rec.Authorizers[authorizer] = true
	}
	rec.Events = append(rec.Events, GrantEvent{
		Timestamp: now,
		Actor:     authorizer,
		Action:    "create",
		Status:    GrantStatusPending,
	})
	r.grants[id] = rec
	return id
}

// Authorize adds an address to the list of permitted signers for the grant.
func (r *GrantRegistry) Authorize(id uint64, actor string) error {
	if actor == "" {
		return errors.New("authorizer required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.grants[id]
	if !ok {
		return errors.New("grant not found")
	}
	if g.Authorizers == nil {
		g.Authorizers = make(map[string]bool)
	}
	if !g.Authorizers[actor] {
		g.Authorizers[actor] = true
		now := time.Now().UTC()
		g.Events = append(g.Events, GrantEvent{
			Timestamp: now,
			Actor:     actor,
			Action:    "authorize",
			Status:    g.Status,
		})
		g.UpdatedAt = now
	}
	return nil
}

// Disburse releases a portion of the grant and records an audit entry.
func (r *GrantRegistry) Disburse(id uint64, amount uint64, note, actor string) error {
	if amount == 0 {
		return errors.New("amount must be positive")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.grants[id]
	if !ok {
		return errors.New("grant not found")
	}
	if len(g.Authorizers) > 0 && !g.Authorizers[actor] {
		return errors.New("actor not authorized")
	}
	if g.Status == GrantStatusRevoked {
		return errors.New("grant revoked")
	}
	if g.Released+amount > g.Amount {
		return errors.New("insufficient remaining funds")
	}
	g.Released += amount
	if note != "" {
		g.Notes = append(g.Notes, note)
	}
	if g.Released == g.Amount {
		g.Status = GrantStatusCompleted
	} else if g.Released > 0 {
		g.Status = GrantStatusActive
	}
	now := time.Now().UTC()
	g.Events = append(g.Events, GrantEvent{
		Timestamp: now,
		Actor:     actor,
		Action:    "disburse",
		Amount:    amount,
		Note:      note,
		Status:    g.Status,
	})
	g.UpdatedAt = now
	return nil
}

// GetGrant returns a grant record by ID.
func (r *GrantRegistry) GetGrant(id uint64) (*GrantRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.grants[id]
	if !ok {
		return nil, false
	}
	return g.clone(), true
}

// ListGrants returns all grants.
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

// Audit returns the audit trail for the specified grant.
func (r *GrantRegistry) Audit(id uint64) ([]GrantEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.grants[id]
	if !ok {
		return nil, errors.New("grant not found")
	}
	events := append([]GrantEvent{}, g.Events...)
	sort.Slice(events, func(i, j int) bool { return events[i].Timestamp.Before(events[j].Timestamp) })
	return events, nil
}

// GrantTelemetry summarises the registry for dashboards and CLI status commands.
type GrantTelemetry struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Active    int `json:"active"`
	Completed int `json:"completed"`
	Revoked   int `json:"revoked"`
}

// Telemetry returns aggregate counters for the registry.
func (r *GrantRegistry) Telemetry() GrantTelemetry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tele := GrantTelemetry{}
	for _, g := range r.grants {
		tele.Total++
		switch g.Status {
		case GrantStatusPending:
			tele.Pending++
		case GrantStatusActive:
			tele.Active++
		case GrantStatusCompleted:
			tele.Completed++
		case GrantStatusRevoked:
			tele.Revoked++
		}
	}
	return tele
}

// Snapshot returns a serialisable copy of all grants for persistence.
func (r *GrantRegistry) Snapshot() []*GrantRecord {
	return r.ListGrants()
}
