package core

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// GrantStatus represents lifecycle states for grants.
type GrantStatus string

const (
	// GrantStatusPending indicates no funds have been released.
	GrantStatusPending GrantStatus = "pending"
	// GrantStatusActive indicates at least one disbursement has occurred.
	GrantStatusActive GrantStatus = "active"
	// GrantStatusCompleted indicates the grant has been fully disbursed.
	GrantStatusCompleted GrantStatus = "completed"
)

// GrantEvent captures audit information for privileged operations.
type GrantEvent struct {
	Timestamp time.Time `json:"Timestamp"`
	Actor     string    `json:"Actor"`
	Action    string    `json:"Action"`
	Amount    uint64    `json:"Amount,omitempty"`
	Note      string    `json:"Note,omitempty"`
}

// GrantRecord captures metadata for a SYN3800 grant token.
type GrantRecord struct {
	ID          uint64              `json:"id"`
	Beneficiary string              `json:"beneficiary"`
	Name        string              `json:"name"`
	Amount      uint64              `json:"amount"`
	Released    uint64              `json:"released"`
	Notes       []string            `json:"notes"`
	Status      GrantStatus         `json:"status"`
	Authorizers map[string]struct{} `json:"authorizers"`
	Audit       []GrantEvent        `json:"audit"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

func (g *GrantRecord) clone() *GrantRecord {
	cp := *g
	if len(g.Notes) > 0 {
		cp.Notes = append([]string(nil), g.Notes...)
	}
	if len(g.Audit) > 0 {
		cp.Audit = append([]GrantEvent(nil), g.Audit...)
	}
	if len(g.Authorizers) > 0 {
		cp.Authorizers = make(map[string]struct{}, len(g.Authorizers))
		for k, v := range g.Authorizers {
			cp.Authorizers[k] = v
		}
	}
	return &cp
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
	id, _ := r.CreateGrantWithAuthorizer(beneficiary, name, amount, "")
	return id
}

// CreateGrantWithAuthorizer registers a new grant linked to an authoriser wallet.
func (r *GrantRegistry) CreateGrantWithAuthorizer(beneficiary, name string, amount uint64, authorizer string) (uint64, error) {
	if beneficiary == "" {
		return 0, fmt.Errorf("beneficiary required")
	}
	if name == "" {
		return 0, fmt.Errorf("name required")
	}
	if amount == 0 {
		return 0, fmt.Errorf("invalid amount")
	}
	now := time.Now().UTC()
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
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if authorizer != "" {
		record.Authorizers[authorizer] = struct{}{}
	}
	record.Audit = append(record.Audit, GrantEvent{Timestamp: now, Actor: authorizer, Action: "create"})
	r.grants[id] = record
	return id, nil
}

// AuthorizeSigner grants disbursement privileges to the provided wallet address.
func (r *GrantRegistry) AuthorizeSigner(id uint64, actor string) error {
	if actor == "" {
		return fmt.Errorf("authorizer address required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	grant, ok := r.grants[id]
	if !ok {
		return errors.New("grant not found")
	}
	grant.Authorizers[actor] = struct{}{}
	grant.UpdatedAt = time.Now().UTC()
	grant.Audit = append(grant.Audit, GrantEvent{Timestamp: grant.UpdatedAt, Actor: actor, Action: "authorize"})
	return nil
}

// Disburse releases a portion of the grant without authorisation enforcement.
func (r *GrantRegistry) Disburse(id uint64, amount uint64, note string) error {
	return r.DisburseWithAuthorizer(id, amount, note, "")
}

// DisburseWithAuthorizer releases funds verifying controller privileges.
func (r *GrantRegistry) DisburseWithAuthorizer(id uint64, amount uint64, note, actor string) error {
	if amount == 0 {
		return fmt.Errorf("invalid disbursement amount")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	grant, ok := r.grants[id]
	if !ok {
		return errors.New("grant not found")
	}
	if len(grant.Authorizers) > 0 {
		if actor == "" {
			return errors.New("unauthorised disbursement")
		}
		if _, ok := grant.Authorizers[actor]; !ok {
			return errors.New("unauthorised disbursement")
		}
	}
	if grant.Released+amount > grant.Amount {
		return errors.New("insufficient remaining funds")
	}
	grant.Released += amount
	if note != "" {
		grant.Notes = append(grant.Notes, note)
	}
	if grant.Released == grant.Amount {
		grant.Status = GrantStatusCompleted
	} else {
		grant.Status = GrantStatusActive
	}
	grant.UpdatedAt = time.Now().UTC()
	grant.Audit = append(grant.Audit, GrantEvent{Timestamp: grant.UpdatedAt, Actor: actor, Action: "disburse", Amount: amount, Note: note})
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

// StatusSummary returns counts grouped by lifecycle state.
func (r *GrantRegistry) StatusSummary() map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	counts := map[string]int{
		string(GrantStatusPending):   0,
		string(GrantStatusActive):    0,
		string(GrantStatusCompleted): 0,
		"total":                      len(r.grants),
	}
	for _, g := range r.grants {
		counts[string(g.Status)]++
	}
	return counts
}

// AuditTrail exposes a copy of the grant's audit log.
func (r *GrantRegistry) AuditTrail(id uint64) ([]GrantEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	grant, ok := r.grants[id]
	if !ok {
		return nil, errors.New("grant not found")
	}
	events := make([]GrantEvent, len(grant.Audit))
	copy(events, grant.Audit)
	return events, nil
}
