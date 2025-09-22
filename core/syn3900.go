package core

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// BenefitStatus represents the lifecycle of a benefit.
type BenefitStatus string

const (
	BenefitStatusPending  BenefitStatus = "PENDING"
	BenefitStatusClaimed  BenefitStatus = "CLAIMED"
	BenefitStatusApproved BenefitStatus = "APPROVED"
)

// BenefitEvent provides an auditable record of state transitions.
type BenefitEvent struct {
	Timestamp time.Time     `json:"timestamp"`
	Actor     string        `json:"actor"`
	Action    string        `json:"action"`
	Status    BenefitStatus `json:"status"`
}

// BenefitRecord holds metadata for a government benefit token issuance.
type BenefitRecord struct {
	ID        uint64          `json:"id"`
	Recipient string          `json:"recipient"`
	Program   string          `json:"program"`
	Amount    uint64          `json:"amount"`
	Claimed   bool            `json:"claimed"`
	Status    BenefitStatus   `json:"status"`
	Approvals map[string]bool `json:"approvals"`
	Events    []BenefitEvent  `json:"events,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (b *BenefitRecord) clone() *BenefitRecord {
	cp := *b
	if len(b.Approvals) > 0 {
		cp.Approvals = make(map[string]bool, len(b.Approvals))
		for k, v := range b.Approvals {
			cp.Approvals[k] = v
		}
	}
	if len(b.Events) > 0 {
		cp.Events = append([]BenefitEvent{}, b.Events...)
	}
	return &cp
}

// BenefitRegistry manages benefit records.
type BenefitRegistry struct {
	mu       sync.RWMutex
	benefits map[uint64]*BenefitRecord
	nextID   uint64
}

// NewBenefitRegistry creates a new registry.
func NewBenefitRegistry() *BenefitRegistry {
	return &BenefitRegistry{benefits: make(map[uint64]*BenefitRecord)}
}

// NewBenefitRegistryFromRecords restores a registry snapshot.
func NewBenefitRegistryFromRecords(records []*BenefitRecord) *BenefitRegistry {
	reg := NewBenefitRegistry()
	var max uint64
	for _, rec := range records {
		if rec == nil {
			continue
		}
		cp := rec.clone()
		reg.benefits[cp.ID] = cp
		if cp.ID > max {
			max = cp.ID
		}
	}
	reg.nextID = max
	return reg
}

// RegisterBenefit records a new benefit and returns its ID.
func (r *BenefitRegistry) RegisterBenefit(recipient, program string, amount uint64, approver string) uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	id := r.nextID
	now := time.Now().UTC()
	rec := &BenefitRecord{
		ID:        id,
		Recipient: recipient,
		Program:   program,
		Amount:    amount,
		Status:    BenefitStatusPending,
		Approvals: make(map[string]bool),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if approver != "" {
		rec.Approvals[approver] = true
	}
	rec.Events = append(rec.Events, BenefitEvent{Timestamp: now, Actor: approver, Action: "register", Status: rec.Status})
	r.benefits[id] = rec
	return id
}

// Claim marks the benefit as claimed by the recipient.
func (r *BenefitRegistry) Claim(id uint64, actor string) error {
	if actor == "" {
		return errors.New("claim actor required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	if b.Recipient != "" && b.Recipient != actor {
		return errors.New("actor not recipient")
	}
	if b.Status == BenefitStatusApproved {
		return nil
	}
	if b.Status == BenefitStatusClaimed {
		return errors.New("benefit already claimed")
	}
	b.Claimed = true
	b.Status = BenefitStatusClaimed
	now := time.Now().UTC()
	b.Events = append(b.Events, BenefitEvent{Timestamp: now, Actor: actor, Action: "claim", Status: b.Status})
	b.UpdatedAt = now
	return nil
}

// Approve marks the benefit as approved by a given actor.
func (r *BenefitRegistry) Approve(id uint64, actor string) error {
	if actor == "" {
		return errors.New("approval actor required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	if b.Approvals == nil {
		b.Approvals = make(map[string]bool)
	}
	if b.Approvals[actor] {
		return nil
	}
	b.Approvals[actor] = true
	b.Status = BenefitStatusApproved
	now := time.Now().UTC()
	b.Events = append(b.Events, BenefitEvent{Timestamp: now, Actor: actor, Action: "approve", Status: b.Status})
	b.UpdatedAt = now
	return nil
}

// GetBenefit retrieves a benefit by ID.
func (r *BenefitRegistry) GetBenefit(id uint64) (*BenefitRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.benefits[id]
	if !ok {
		return nil, false
	}
	return b.clone(), true
}

// ListBenefits returns all benefit records.
func (r *BenefitRegistry) ListBenefits() []*BenefitRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*BenefitRecord, 0, len(r.benefits))
	for _, b := range r.benefits {
		res = append(res, b.clone())
	}
	sort.Slice(res, func(i, j int) bool { return res[i].ID < res[j].ID })
	return res
}

// BenefitTelemetry summarises benefit activity for dashboards.
type BenefitTelemetry struct {
	Total    int `json:"total"`
	Pending  int `json:"pending"`
	Claimed  int `json:"claimed"`
	Approved int `json:"approved"`
}

// Telemetry reports aggregate counters.
func (r *BenefitRegistry) Telemetry() BenefitTelemetry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tele := BenefitTelemetry{}
	for _, b := range r.benefits {
		tele.Total++
		switch b.Status {
		case BenefitStatusPending:
			tele.Pending++
		case BenefitStatusClaimed:
			tele.Claimed++
		case BenefitStatusApproved:
			tele.Approved++
		}
	}
	return tele
}

// Snapshot returns a serialisable copy of the registry.
func (r *BenefitRegistry) Snapshot() []*BenefitRecord {
	return r.ListBenefits()
}
