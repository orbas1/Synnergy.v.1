package tokens

import (
	"errors"
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
	Roles       map[string]bool
	Status      string
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
		Roles:       make(map[string]bool),
		Status:      "active",
	}
}

// UpdateVotingPower sets the member's voting power.
func (m *Syn2500Member) UpdateVotingPower(power uint64) {
	m.VotingPower = power
}

// AssignRole grants a governance role to the member.
func (m *Syn2500Member) AssignRole(role string) {
	if m.Roles == nil {
		m.Roles = make(map[string]bool)
	}
	m.Roles[role] = true
}

// RemoveRole revokes a governance role.
func (m *Syn2500Member) RemoveRole(role string) {
	delete(m.Roles, role)
}

// SetStatus updates the participation status of the member.
func (m *Syn2500Member) SetStatus(status string) {
	m.Status = status
}

// Syn2500Registry manages DAO membership records.
type Syn2500Registry struct {
	mu      sync.RWMutex
	members map[string]*Syn2500Member
	events  []Syn2500Event
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
	r.events = append(r.events, Syn2500Event{Type: "add", MemberID: m.ID, Timestamp: time.Now(), Metadata: map[string]string{"address": m.Address}})
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
	r.events = append(r.events, Syn2500Event{Type: "remove", MemberID: id, Timestamp: time.Now()})
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

// UpdateMetadata updates metadata for a member.
func (r *Syn2500Registry) UpdateMetadata(id string, meta map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	member, ok := r.members[id]
	if !ok {
		return errors.New("member not found")
	}
	if member.Metadata == nil {
		member.Metadata = make(map[string]string)
	}
	for k, v := range meta {
		member.Metadata[k] = v
	}
	r.events = append(r.events, Syn2500Event{Type: "metadata", MemberID: id, Timestamp: time.Now(), Metadata: meta})
	return nil
}

// TotalVotingPower returns the aggregate voting power of active members.
func (r *Syn2500Registry) TotalVotingPower() uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var total uint64
	for _, m := range r.members {
		if m.Status == "active" {
			total += m.VotingPower
		}
	}
	return total
}

// MembersWithRole filters members by assigned role.
func (r *Syn2500Registry) MembersWithRole(role string) []*Syn2500Member {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Syn2500Member, 0)
	for _, m := range r.members {
		if m.Roles[role] {
			out = append(out, m)
		}
	}
	return out
}

// Syn2500Event captures registry lifecycle actions for audit trails.
type Syn2500Event struct {
	Type      string
	MemberID  string
	Timestamp time.Time
	Metadata  map[string]string
}

// Events returns tracked registry actions.
func (r *Syn2500Registry) Events(limit int) []Syn2500Event {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if limit <= 0 || limit >= len(r.events) {
		out := make([]Syn2500Event, len(r.events))
		copy(out, r.events)
		return out
	}
	out := make([]Syn2500Event, limit)
	copy(out, r.events[len(r.events)-limit:])
	return out
}
