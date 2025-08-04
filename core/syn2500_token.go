package core

import (
	"sync"
	"time"
)

// Syn2500Member stores metadata about DAO membership.
type Syn2500Member struct {
	ID          string
	Address     string
	JoinedAt    time.Time
	VotingPower uint64
	Metadata    map[string]string
}

// NewSyn2500Member creates a new DAO member record.
func NewSyn2500Member(id, addr string, power uint64, meta map[string]string) *Syn2500Member {
	cp := make(map[string]string, len(meta))
	for k, v := range meta {
		cp[k] = v
	}
	return &Syn2500Member{
		ID:          id,
		Address:     addr,
		JoinedAt:    time.Now(),
		VotingPower: power,
		Metadata:    cp,
	}
}

// UpdateVotingPower sets the member's voting power.
func (m *Syn2500Member) UpdateVotingPower(power uint64) {
	m.VotingPower = power
}

// Syn2500Registry manages DAO membership records.
type Syn2500Registry struct {
	mu      sync.RWMutex
	members map[string]*Syn2500Member
}

// NewSyn2500Registry returns an empty membership registry.
func NewSyn2500Registry() *Syn2500Registry {
	return &Syn2500Registry{members: make(map[string]*Syn2500Member)}
}

// AddMember inserts or replaces a member entry.
func (r *Syn2500Registry) AddMember(m *Syn2500Member) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.members[m.ID] = m
}

// GetMember retrieves a member by ID.
func (r *Syn2500Registry) GetMember(id string) (*Syn2500Member, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.members[id]
	return m, ok
}

// RemoveMember deletes a member from the registry.
func (r *Syn2500Registry) RemoveMember(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.members, id)
}

// ListMembers returns all members of the registry.
func (r *Syn2500Registry) ListMembers() []*Syn2500Member {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*Syn2500Member, 0, len(r.members))
	for _, m := range r.members {
		list = append(list, m)
	}
	return list
}
