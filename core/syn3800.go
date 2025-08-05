package core

import (
	"errors"
	"sync"
)

// GrantRecord captures metadata for a SYN3800 grant token.
type GrantRecord struct {
	ID          uint64
	Beneficiary string
	Name        string
	Amount      uint64
	Released    uint64
	Notes       []string
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
	r.grants[id] = &GrantRecord{ID: id, Beneficiary: beneficiary, Name: name, Amount: amount}
	return id
}

// Disburse releases a portion of the grant.
func (r *GrantRegistry) Disburse(id uint64, amount uint64, note string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.grants[id]
	if !ok {
		return errors.New("grant not found")
	}
	if g.Released+amount > g.Amount {
		return errors.New("insufficient remaining funds")
	}
	g.Released += amount
	if note != "" {
		g.Notes = append(g.Notes, note)
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
	return &cp, true
}

// ListGrants returns all grants.
func (r *GrantRegistry) ListGrants() []*GrantRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*GrantRecord, 0, len(r.grants))
	for _, g := range r.grants {
		cp := *g
		res = append(res, &cp)
	}
	return res
}
