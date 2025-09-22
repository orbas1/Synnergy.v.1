package core

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// BenefitStatus describes a benefit's lifecycle stage.
type BenefitStatus string

const (
	BenefitStatusPending  BenefitStatus = "pending"
	BenefitStatusClaimed  BenefitStatus = "claimed"
	BenefitStatusApproved BenefitStatus = "approved"
)

// BenefitEvent records auditable actions against a benefit enrolment.
type BenefitEvent struct {
	Timestamp time.Time     `json:"timestamp"`
	Type      string        `json:"type"`
	Actor     string        `json:"actor,omitempty"`
	Note      string        `json:"note,omitempty"`
	Status    BenefitStatus `json:"status"`
}

// BenefitRecord holds metadata for a government benefit token issuance.
type BenefitRecord struct {
	ID        uint64         `json:"id"`
	Recipient string         `json:"recipient"`
	Program   string         `json:"program"`
	Amount    uint64         `json:"amount"`
	Claimed   bool           `json:"claimed"`
	Status    BenefitStatus  `json:"status"`
	Approvers []string       `json:"approvers"`
	Events    []BenefitEvent `json:"events"`
}

type benefitState struct {
	record    BenefitRecord
	approvers map[string]struct{}
}

// BenefitRegistry manages benefit records.
type BenefitRegistry struct {
	mu       sync.RWMutex
	benefits map[uint64]*benefitState
	nextID   uint64
}

// NewBenefitRegistry creates a new registry.
func NewBenefitRegistry() *BenefitRegistry {
	return &BenefitRegistry{benefits: make(map[uint64]*benefitState)}
}

// RegisterBenefit records a new benefit and returns its ID.
func (r *BenefitRegistry) RegisterBenefit(recipient, program string, amount uint64) uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	id := r.nextID
	state := &benefitState{
		record: BenefitRecord{
			ID:        id,
			Recipient: recipient,
			Program:   program,
			Amount:    amount,
			Status:    BenefitStatusPending,
			Events: []BenefitEvent{{
				Timestamp: time.Now().UTC(),
				Type:      "registered",
				Actor:     recipient,
				Status:    BenefitStatusPending,
			}},
		},
		approvers: make(map[string]struct{}),
	}
	r.benefits[id] = state
	return id
}

// AddApprover registers a wallet address as an authorised approver.
func (r *BenefitRegistry) AddApprover(id uint64, addr string) error {
	if addr == "" {
		return errors.New("approver required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	state, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	if _, exists := state.approvers[addr]; exists {
		return nil
	}
	state.approvers[addr] = struct{}{}
	state.record.Approvers = append(state.record.Approvers, addr)
	state.record.Events = append(state.record.Events, BenefitEvent{
		Timestamp: time.Now().UTC(),
		Type:      "approver_added",
		Actor:     addr,
		Status:    state.record.Status,
	})
	sort.Strings(state.record.Approvers)
	return nil
}

// Claim marks the benefit as claimed by the recipient wallet.
func (r *BenefitRegistry) Claim(id uint64, actor string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	state, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	if actor == "" || actor != state.record.Recipient {
		return errors.New("authentication failed")
	}
	state.record.Claimed = true
	state.record.Status = BenefitStatusClaimed
	state.record.Events = append(state.record.Events, BenefitEvent{
		Timestamp: time.Now().UTC(),
		Type:      "claimed",
		Actor:     actor,
		Status:    state.record.Status,
	})
	return nil
}

// Approve records an approval from a trusted wallet.
func (r *BenefitRegistry) Approve(id uint64, actor string) error {
	if actor == "" {
		return errors.New("authentication failed")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	state, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	if len(state.approvers) > 0 {
		if _, authorised := state.approvers[actor]; !authorised {
			return errors.New("approver not authorised")
		}
	}
	state.record.Status = BenefitStatusApproved
	if !state.record.Claimed {
		state.record.Claimed = true
	}
	state.record.Events = append(state.record.Events, BenefitEvent{
		Timestamp: time.Now().UTC(),
		Type:      "approved",
		Actor:     actor,
		Status:    state.record.Status,
	})
	return nil
}

// GetBenefit retrieves a benefit by ID.
func (r *BenefitRegistry) GetBenefit(id uint64) (*BenefitRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	state, ok := r.benefits[id]
	if !ok {
		return nil, false
	}
	cp := state.record
	cp.Approvers = append([]string(nil), cp.Approvers...)
	cp.Events = append([]BenefitEvent(nil), cp.Events...)
	return &cp, true
}

// ListBenefits returns a copy of all records.
func (r *BenefitRegistry) ListBenefits() []*BenefitRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*BenefitRecord, 0, len(r.benefits))
	for _, state := range r.benefits {
		cp := state.record
		cp.Approvers = append([]string(nil), cp.Approvers...)
		cp.Events = append([]BenefitEvent(nil), cp.Events...)
		res = append(res, &cp)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].ID < res[j].ID })
	return res
}

// BenefitSummary provides telemetry statistics for benefits.
type BenefitSummary struct {
	Total    int `json:"total"`
	Pending  int `json:"pending"`
	Claimed  int `json:"claimed"`
	Approved int `json:"approved"`
}

// Summary aggregates status counts.
func (r *BenefitRegistry) Summary() BenefitSummary {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var summary BenefitSummary
	summary.Total = len(r.benefits)
	for _, state := range r.benefits {
		switch state.record.Status {
		case BenefitStatusPending:
			summary.Pending++
		case BenefitStatusClaimed:
			summary.Claimed++
		case BenefitStatusApproved:
			summary.Approved++
			summary.Claimed++
		}
	}
	return summary
}

// Events returns a copy of the event stream for the benefit.
func (r *BenefitRegistry) Events(id uint64) ([]BenefitEvent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	state, ok := r.benefits[id]
	if !ok {
		return nil, false
	}
	ev := append([]BenefitEvent(nil), state.record.Events...)
	return ev, true
}
