package core

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// GrantStatus describes the lifecycle state for a SYN3800 grant record.
type GrantStatus string

const (
	// GrantStatusPending represents a newly created grant with no releases.
	GrantStatusPending GrantStatus = "pending"
	// GrantStatusActive indicates at least one authorised release occurred.
	GrantStatusActive GrantStatus = "active"
	// GrantStatusCompleted indicates the full amount has been released.
	GrantStatusCompleted GrantStatus = "completed"
)

// GrantEvent captures an auditable change to a grant record.
type GrantEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Amount    uint64    `json:"amount,omitempty"`
	Note      string    `json:"note,omitempty"`
	Actor     string    `json:"actor,omitempty"`
}

// GrantRecord captures metadata for a SYN3800 grant token.
type GrantRecord struct {
	ID          uint64       `json:"id"`
	Beneficiary string       `json:"beneficiary"`
	Name        string       `json:"name"`
	Amount      uint64       `json:"amount"`
	Released    uint64       `json:"released"`
	Notes       []string     `json:"notes,omitempty"`
	Status      GrantStatus  `json:"status"`
	Authorizers []string     `json:"authorizers"`
	Events      []GrantEvent `json:"events"`
}

type grantState struct {
	record      GrantRecord
	authorizers map[string]struct{}
}

// GrantRegistry manages grant records.
type GrantRegistry struct {
	mu     sync.RWMutex
	grants map[uint64]*grantState
	nextID uint64
}

// NewGrantRegistry creates a new registry.
func NewGrantRegistry() *GrantRegistry {
	return &GrantRegistry{grants: make(map[uint64]*grantState)}
}

// CreateGrant registers a new grant and returns its ID.
func (r *GrantRegistry) CreateGrant(beneficiary, name string, amount uint64) uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	id := r.nextID
	state := &grantState{
		record: GrantRecord{
			ID:          id,
			Beneficiary: beneficiary,
			Name:        name,
			Amount:      amount,
			Status:      GrantStatusPending,
			Events: []GrantEvent{{
				Timestamp: time.Now().UTC(),
				Type:      "created",
				Actor:     beneficiary,
				Amount:    amount,
			}},
		},
		authorizers: make(map[string]struct{}),
	}
	r.grants[id] = state
	return id
}

// AddAuthorizer registers a wallet address as an authorised releaser for the grant.
func (r *GrantRegistry) AddAuthorizer(id uint64, address string) error {
	if address == "" {
		return errors.New("authorizer required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	state, ok := r.grants[id]
	if !ok {
		return errors.New("grant not found")
	}
	if _, exists := state.authorizers[address]; exists {
		return nil
	}
	state.authorizers[address] = struct{}{}
	state.record.Authorizers = append(state.record.Authorizers, address)
	state.record.Events = append(state.record.Events, GrantEvent{
		Timestamp: time.Now().UTC(),
		Type:      "authorized",
		Actor:     address,
	})
	sort.Strings(state.record.Authorizers)
	return nil
}

// Disburse releases a portion of the grant without enforcing authorisation.
// Legacy callers without wallet context can continue to use this helper.
func (r *GrantRegistry) Disburse(id uint64, amount uint64, note string) error {
	return r.DisburseWithActor(id, amount, note, "")
}

// DisburseWithActor releases a portion of the grant ensuring the actor is authorised.
func (r *GrantRegistry) DisburseWithActor(id uint64, amount uint64, note, actor string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	state, ok := r.grants[id]
	if !ok {
		return errors.New("grant not found")
	}
	if amount == 0 {
		return errors.New("amount must be positive")
	}
	if len(state.authorizers) > 0 {
		if actor == "" {
			return errors.New("authorizer required")
		}
		if _, authorised := state.authorizers[actor]; !authorised {
			return fmt.Errorf("unauthorised actor %s", actor)
		}
	}
	if state.record.Released+amount > state.record.Amount {
		return errors.New("insufficient remaining funds")
	}
	state.record.Released += amount
	if note != "" {
		state.record.Notes = append(state.record.Notes, note)
	}
	switch {
	case state.record.Released == state.record.Amount:
		state.record.Status = GrantStatusCompleted
	case state.record.Released > 0:
		state.record.Status = GrantStatusActive
	default:
		state.record.Status = GrantStatusPending
	}
	state.record.Events = append(state.record.Events, GrantEvent{
		Timestamp: time.Now().UTC(),
		Type:      "disbursed",
		Amount:    amount,
		Note:      note,
		Actor:     actor,
	})
	return nil
}

// GetGrant returns a grant record by ID.
func (r *GrantRegistry) GetGrant(id uint64) (*GrantRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	state, ok := r.grants[id]
	if !ok {
		return nil, false
	}
	cp := state.record
	cp.Authorizers = append([]string(nil), cp.Authorizers...)
	cp.Events = append([]GrantEvent(nil), cp.Events...)
	cp.Notes = append([]string(nil), cp.Notes...)
	return &cp, true
}

// ListGrants returns all grants.
func (r *GrantRegistry) ListGrants() []*GrantRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*GrantRecord, 0, len(r.grants))
	for _, state := range r.grants {
		cp := state.record
		cp.Authorizers = append([]string(nil), cp.Authorizers...)
		cp.Events = append([]GrantEvent(nil), cp.Events...)
		cp.Notes = append([]string(nil), cp.Notes...)
		res = append(res, &cp)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].ID < res[j].ID })
	return res
}

// GrantSummary aggregates totals for telemetry reporting.
type GrantSummary struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Active    int `json:"active"`
	Completed int `json:"completed"`
}

// Summary computes aggregate status counts.
func (r *GrantRegistry) Summary() GrantSummary {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var summary GrantSummary
	summary.Total = len(r.grants)
	for _, state := range r.grants {
		switch state.record.Status {
		case GrantStatusPending:
			summary.Pending++
		case GrantStatusActive:
			summary.Active++
		case GrantStatusCompleted:
			summary.Completed++
		}
	}
	return summary
}

// AuditTrail returns the recorded events for a grant.
func (r *GrantRegistry) AuditTrail(id uint64) ([]GrantEvent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	state, ok := r.grants[id]
	if !ok {
		return nil, false
	}
	events := append([]GrantEvent(nil), state.record.Events...)
	return events, true
}
