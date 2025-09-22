package core

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// BenefitStatus represents the lifecycle state for SYN3900 benefits.
type BenefitStatus string

const (
	BenefitStatusPending  BenefitStatus = "pending"
	BenefitStatusClaimed  BenefitStatus = "claimed"
	BenefitStatusApproved BenefitStatus = "approved"
)

// BenefitRecord holds metadata for a government benefit token issuance.
type BenefitRecord struct {
	ID        uint64              `json:"id"`
	Recipient string              `json:"recipient"`
	Program   string              `json:"program"`
	Amount    uint64              `json:"amount"`
	Claimed   bool                `json:"claimed"`
	Status    BenefitStatus       `json:"status"`
	Approvals map[string]struct{} `json:"-"`
	Claimant  string              `json:"claimant"`
	CreatedAt time.Time           `json:"created_at"`
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

// RegisterBenefit records a new benefit and returns its ID.
func (r *BenefitRegistry) RegisterBenefit(recipient, program string, amount uint64) uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	id := r.nextID
	r.benefits[id] = &BenefitRecord{
		ID:        id,
		Recipient: recipient,
		Program:   program,
		Amount:    amount,
		Status:    BenefitStatusPending,
		Approvals: make(map[string]struct{}),
		CreatedAt: time.Now().UTC(),
	}
	return id
}

// Claim marks the benefit as claimed by signer.
func (r *BenefitRegistry) Claim(id uint64, signer string) error {
	if signer == "" {
		return errors.New("claim signer required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	if b.Claimed {
		return errors.New("benefit already claimed")
	}
	if b.Recipient != "" && b.Recipient != signer {
		return errors.New("claimant must match recipient")
	}
	b.Claimed = true
	b.Claimant = signer
	b.Status = BenefitStatusClaimed
	return nil
}

// Approve records an approval for the benefit.
func (r *BenefitRegistry) Approve(id uint64, signer string) error {
	if signer == "" {
		return errors.New("approver required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	if b.Approvals == nil {
		b.Approvals = make(map[string]struct{})
	}
	b.Approvals[signer] = struct{}{}
	b.Status = BenefitStatusApproved
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
	cp := copyBenefit(b)
	return &cp, true
}

// ListBenefits returns all benefit records.
func (r *BenefitRegistry) ListBenefits() []*BenefitRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*BenefitRecord, 0, len(r.benefits))
	ids := make([]uint64, 0, len(r.benefits))
	for id := range r.benefits {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	for _, id := range ids {
		cp := copyBenefit(r.benefits[id])
		res = append(res, &cp)
	}
	return res
}

// BenefitTelemetry summarises registry metrics.
type BenefitTelemetry struct {
	Total    int `json:"total"`
	Pending  int `json:"pending"`
	Approved int `json:"approved"`
	Claimed  int `json:"claimed"`
}

// Telemetry aggregates status counts for monitoring.
func (r *BenefitRegistry) Telemetry() BenefitTelemetry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tele := BenefitTelemetry{}
	for _, b := range r.benefits {
		tele.Total++
		switch b.Status {
		case BenefitStatusApproved:
			tele.Approved++
			tele.Claimed++
		case BenefitStatusClaimed:
			tele.Claimed++
		default:
			tele.Pending++
		}
	}
	return tele
}

func copyBenefit(b *BenefitRecord) BenefitRecord {
	cp := BenefitRecord{
		ID:        b.ID,
		Recipient: b.Recipient,
		Program:   b.Program,
		Amount:    b.Amount,
		Claimed:   b.Claimed,
		Status:    b.Status,
		Claimant:  b.Claimant,
		CreatedAt: b.CreatedAt,
	}
	if len(b.Approvals) > 0 {
		cp.Approvals = make(map[string]struct{}, len(b.Approvals))
		for addr := range b.Approvals {
			cp.Approvals[addr] = struct{}{}
		}
	}
	return cp
}
