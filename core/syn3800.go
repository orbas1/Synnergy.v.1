package core

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// GrantStatus represents the lifecycle state of a SYN3800 grant.
type GrantStatus string

const (
	GrantStatusPending   GrantStatus = "pending"
	GrantStatusActive    GrantStatus = "active"
	GrantStatusCompleted GrantStatus = "completed"
)

// GrantEventType captures the type of lifecycle mutation for a grant.
type GrantEventType string

const (
	GrantEventCreated    GrantEventType = "created"
	GrantEventAuthorized GrantEventType = "authorized"
	GrantEventDisbursed  GrantEventType = "disbursed"
)

// GrantEvent records an auditable action applied to a grant.
type GrantEvent struct {
	Timestamp time.Time      `json:"timestamp"`
	GrantID   uint64         `json:"grant_id"`
	Type      GrantEventType `json:"type"`
	Amount    uint64         `json:"amount,omitempty"`
	Note      string         `json:"note,omitempty"`
	Signer    string         `json:"signer,omitempty"`
}

// GrantRecord captures metadata for a SYN3800 grant token.
type GrantRecord struct {
	ID          uint64              `json:"id"`
	Beneficiary string              `json:"beneficiary"`
	Name        string              `json:"name"`
	Amount      uint64              `json:"amount"`
	Released    uint64              `json:"released"`
	Status      GrantStatus         `json:"status"`
	Authorizers map[string]struct{} `json:"-"`
	Events      []GrantEvent        `json:"events"`
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

// CreateGrant registers a new grant and returns its ID.
func (r *GrantRegistry) CreateGrant(beneficiary, name string, amount uint64) uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	id := r.nextID
	record := &GrantRecord{
		ID:          id,
		Beneficiary: beneficiary,
		Name:        name,
		Amount:      amount,
		Status:      GrantStatusPending,
		Authorizers: make(map[string]struct{}),
		Events: []GrantEvent{{
			Timestamp: time.Now().UTC(),
			GrantID:   id,
			Type:      GrantEventCreated,
		}},
	}
	r.grants[id] = record
	return id
}

// Authorize registers signer as an approved disbursement authority for the grant.
func (r *GrantRegistry) Authorize(id uint64, signer string) error {
	if signer == "" {
		return errors.New("authorizer address required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.grants[id]
	if !ok {
		return errors.New("grant not found")
	}
	if g.Authorizers == nil {
		g.Authorizers = make(map[string]struct{})
	}
	if _, exists := g.Authorizers[signer]; exists {
		return nil
	}
	g.Authorizers[signer] = struct{}{}
	g.Events = append(g.Events, GrantEvent{
		Timestamp: time.Now().UTC(),
		GrantID:   id,
		Type:      GrantEventAuthorized,
		Signer:    signer,
	})
	return nil
}

// Disburse releases a portion of the grant using signer authority.
func (r *GrantRegistry) Disburse(id uint64, amount uint64, note, signer string) error {
	if amount == 0 {
		return fmt.Errorf("amount must be positive")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.grants[id]
	if !ok {
		return errors.New("grant not found")
	}
	if len(g.Authorizers) > 0 {
		if _, authorized := g.Authorizers[signer]; !authorized {
			return errors.New("authorizer required")
		}
	}
	if g.Released+amount > g.Amount {
		return errors.New("insufficient remaining funds")
	}
	g.Released += amount
	if g.Released == g.Amount {
		g.Status = GrantStatusCompleted
	} else {
		g.Status = GrantStatusActive
	}
	g.Events = append(g.Events, GrantEvent{
		Timestamp: time.Now().UTC(),
		GrantID:   id,
		Type:      GrantEventDisbursed,
		Amount:    amount,
		Note:      note,
		Signer:    signer,
	})
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
	cp := copyGrant(g)
	return &cp, true
}

// ListGrants returns all grants.
func (r *GrantRegistry) ListGrants() []*GrantRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*GrantRecord, 0, len(r.grants))
	ids := make([]uint64, 0, len(r.grants))
	for id := range r.grants {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	for _, id := range ids {
		cp := copyGrant(r.grants[id])
		res = append(res, &cp)
	}
	return res
}

// Audit returns the immutable event log for the grant.
func (r *GrantRegistry) Audit(id uint64) ([]GrantEvent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.grants[id]
	if !ok {
		return nil, false
	}
	events := make([]GrantEvent, len(g.Events))
	copy(events, g.Events)
	return events, true
}

// GrantTelemetry summarises registry activity for monitoring endpoints.
type GrantTelemetry struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Active    int `json:"active"`
	Completed int `json:"completed"`
}

// Telemetry reports aggregate lifecycle counts across grants.
func (r *GrantRegistry) Telemetry() GrantTelemetry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tele := GrantTelemetry{}
	for _, g := range r.grants {
		tele.Total++
		switch g.Status {
		case GrantStatusCompleted:
			tele.Completed++
		case GrantStatusActive:
			tele.Active++
		default:
			tele.Pending++
		}
	}
	return tele
}

func copyGrant(g *GrantRecord) GrantRecord {
	cp := GrantRecord{
		ID:          g.ID,
		Beneficiary: g.Beneficiary,
		Name:        g.Name,
		Amount:      g.Amount,
		Released:    g.Released,
		Status:      g.Status,
	}
	if len(g.Authorizers) > 0 {
		cp.Authorizers = make(map[string]struct{}, len(g.Authorizers))
		for k := range g.Authorizers {
			cp.Authorizers[k] = struct{}{}
		}
	}
	if len(g.Events) > 0 {
		cp.Events = make([]GrantEvent, len(g.Events))
		copy(cp.Events, g.Events)
	}
	return cp
}
