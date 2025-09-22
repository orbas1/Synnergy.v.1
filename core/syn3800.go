package core

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// GrantStatus captures the lifecycle stage of a SYN3800 grant.
type GrantStatus string

const (
	GrantStatusDraft      GrantStatus = "draft"
	GrantStatusPending    GrantStatus = "pending"
	GrantStatusAuthorized GrantStatus = "authorized"
	GrantStatusReleased   GrantStatus = "released"
	GrantStatusActive     GrantStatus = GrantStatusAuthorized
	GrantStatusCompleted  GrantStatus = GrantStatusReleased
)

// GrantEvent maintains backwards compatibility with existing JSON consumers.
type GrantEvent = GrantAuditEvent

// GrantAuditEvent records governance actions for transparency.
type GrantAuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Actor     Address   `json:"actor"`
	Action    string    `json:"action"`
	Amount    uint64    `json:"amount,omitempty"`
	Note      string    `json:"note,omitempty"`
}

// GrantRecord captures metadata for a SYN3800 grant token.
type GrantRecord struct {
	ID          uint64                `json:"id"`
	Beneficiary string                `json:"beneficiary"`
	Name        string                `json:"name"`
	Amount      uint64                `json:"amount"`
	Released    uint64                `json:"released"`
	Status      GrantStatus           `json:"status"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	Authorizers map[Address]time.Time `json:"authorizers"`
	Audit       []GrantAuditEvent     `json:"audit"`
}

// GrantRegistrySnapshot stores registry state for Stage 73 persistence.
type GrantRegistrySnapshot struct {
	NextID      uint64        `json:"next_id"`
	Grants      []GrantRecord `json:"grants"`
	GeneratedAt time.Time     `json:"generated_at"`
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
func (r *GrantRegistry) CreateGrant(beneficiary, name string, amount uint64, authorizer Address) (uint64, error) {
	if beneficiary == "" || name == "" {
		return 0, fmt.Errorf("beneficiary and name required")
	}
	if amount == 0 {
		return 0, fmt.Errorf("amount must be positive")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	id := r.nextID
	now := time.Now().UTC()
	record := &GrantRecord{
		ID:          id,
		Beneficiary: beneficiary,
		Name:        name,
		Amount:      amount,
		Status:      GrantStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
		Authorizers: make(map[Address]time.Time),
	}
	if authorizer != "" {
		record.Authorizers[authorizer] = now
		record.Status = GrantStatusAuthorized
	}
	record.Audit = append(record.Audit, GrantAuditEvent{Timestamp: now, Actor: authorizer, Action: "create", Amount: amount})
	r.grants[id] = record
	return id, nil
}

// Authorize grants release permissions to a wallet.
func (r *GrantRegistry) Authorize(id uint64, actor Address) error {
	if actor == "" {
		return fmt.Errorf("authorizer required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.grants[id]
	if !ok {
		return errors.New("grant not found")
	}
	if g.Authorizers == nil {
		g.Authorizers = make(map[Address]time.Time)
	}
	if _, exists := g.Authorizers[actor]; exists {
		return nil
	}
	now := time.Now().UTC()
	g.Authorizers[actor] = now
	g.Status = GrantStatusAuthorized
	g.UpdatedAt = now
	g.Audit = append(g.Audit, GrantAuditEvent{Timestamp: now, Actor: actor, Action: "authorize"})
	return nil
}

// Disburse releases a portion of the grant.
func (r *GrantRegistry) Disburse(id uint64, amount uint64, note string, actor Address) error {
	if amount == 0 {
		return fmt.Errorf("amount must be positive")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.grants[id]
	if !ok {
		return errors.New("grant not found")
	}
	if g.Authorizers == nil || actor == "" {
		return fmt.Errorf("authorizer required")
	}
	if _, ok := g.Authorizers[actor]; !ok {
		return fmt.Errorf("wallet not authorised")
	}
	if g.Released+amount > g.Amount {
		return fmt.Errorf("insufficient remaining funds")
	}
	g.Released += amount
	now := time.Now().UTC()
	g.UpdatedAt = now
	g.Audit = append(g.Audit, GrantAuditEvent{Timestamp: now, Actor: actor, Action: "disburse", Amount: amount, Note: note})
	if g.Released == g.Amount {
		g.Status = GrantStatusReleased
	}
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
	cp := *g
	cp.Authorizers = copyAuthorizers(g.Authorizers)
	cp.Audit = append([]GrantAuditEvent(nil), g.Audit...)
	return &cp, true
}

func copyAuthorizers(src map[Address]time.Time) map[Address]time.Time {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[Address]time.Time, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// ListGrants returns all grants.
func (r *GrantRegistry) ListGrants() []GrantRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]GrantRecord, 0, len(r.grants))
	for _, g := range r.grants {
		cp := *g
		cp.Authorizers = copyAuthorizers(g.Authorizers)
		cp.Audit = append([]GrantAuditEvent(nil), g.Audit...)
		res = append(res, cp)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].ID < res[j].ID })
	return res
}

// AuditLog returns the audit events for a grant.
func (r *GrantRegistry) AuditLog(id uint64) []GrantAuditEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.grants[id]
	if !ok {
		return nil
	}
	return append([]GrantAuditEvent(nil), g.Audit...)
}

// GrantRegistryStatus summarises registry health.
type GrantRegistryStatus struct {
	Total       int    `json:"total"`
	Pending     int    `json:"pending"`
	Active      int    `json:"active"`
	Completed   int    `json:"completed"`
	Outstanding uint64 `json:"outstanding"`
}

// StatusSummary aggregates registry metrics for dashboards.
func (r *GrantRegistry) StatusSummary() GrantRegistryStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()
	status := GrantRegistryStatus{}
	for _, g := range r.grants {
		status.Total++
		switch g.Status {
		case GrantStatusPending:
			status.Pending++
		case GrantStatusAuthorized:
			status.Active++
		case GrantStatusReleased:
			status.Completed++
		}
		status.Outstanding += g.Amount - g.Released
	}
	return status
}

// Snapshot captures the registry for persistence.
func (r *GrantRegistry) Snapshot() GrantRegistrySnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	grants := make([]GrantRecord, 0, len(r.grants))
	for _, g := range r.grants {
		cp := *g
		cp.Authorizers = copyAuthorizers(g.Authorizers)
		cp.Audit = append([]GrantAuditEvent(nil), g.Audit...)
		grants = append(grants, cp)
	}
	sort.Slice(grants, func(i, j int) bool { return grants[i].ID < grants[j].ID })
	return GrantRegistrySnapshot{NextID: r.nextID, Grants: grants, GeneratedAt: time.Now().UTC()}
}

// Restore replaces the registry contents with a snapshot.
func (r *GrantRegistry) Restore(snapshot GrantRegistrySnapshot) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID = snapshot.NextID
	r.grants = make(map[uint64]*GrantRecord, len(snapshot.Grants))
	for i := range snapshot.Grants {
		grant := snapshot.Grants[i]
		cp := grant
		cp.Authorizers = copyAuthorizers(grant.Authorizers)
		cp.Audit = append([]GrantAuditEvent(nil), grant.Audit...)
		r.grants[grant.ID] = &cp
	}
}
