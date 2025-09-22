package core

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// BenefitStatus captures government benefit lifecycle.
type BenefitStatus string

const (
	BenefitStatusPending  BenefitStatus = "pending"
	BenefitStatusApproved BenefitStatus = "approved"
	BenefitStatusClaimed  BenefitStatus = "claimed"
)

// BenefitEvent captures audit entries for benefit operations.
type BenefitEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
}

// BenefitRecord holds metadata for a government benefit token issuance.
type BenefitRecord struct {
	ID        uint64              `json:"id"`
	Recipient string              `json:"recipient"`
	Program   string              `json:"program"`
	Amount    uint64              `json:"amount"`
	Claimed   bool                `json:"claimed"`
	Status    BenefitStatus       `json:"status"`
	Approvers map[string]struct{} `json:"approvers"`
	Audit     []BenefitEvent      `json:"audit"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

func (b *BenefitRecord) clone() *BenefitRecord {
	cp := *b
	if len(b.Approvers) > 0 {
		cp.Approvers = make(map[string]struct{}, len(b.Approvers))
		for k, v := range b.Approvers {
			cp.Approvers[k] = v
		}
	}
	if len(b.Audit) > 0 {
		cp.Audit = append([]BenefitEvent(nil), b.Audit...)
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

// RegisterBenefit records a new benefit and returns its ID.
func (r *BenefitRegistry) RegisterBenefit(recipient, program string, amount uint64) uint64 {
	id, _ := r.RegisterBenefitWithApprover(recipient, program, amount, "")
	return id
}

// RegisterBenefitWithApprover records a benefit and initial approval wallet.
func (r *BenefitRegistry) RegisterBenefitWithApprover(recipient, program string, amount uint64, approver string) (uint64, error) {
	if recipient == "" {
		return 0, fmt.Errorf("recipient required")
	}
	if program == "" {
		return 0, fmt.Errorf("program required")
	}
	if amount == 0 {
		return 0, fmt.Errorf("invalid amount")
	}
	now := time.Now().UTC()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	id := r.nextID
	record := &BenefitRecord{
		ID:        id,
		Recipient: recipient,
		Program:   program,
		Amount:    amount,
		Status:    BenefitStatusPending,
		Approvers: make(map[string]struct{}),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if approver != "" {
		record.Approvers[approver] = struct{}{}
		record.Status = BenefitStatusApproved
	}
	record.Audit = append(record.Audit, BenefitEvent{Timestamp: now, Actor: approver, Action: "register"})
	r.benefits[id] = record
	return id, nil
}

// Approve registers an additional wallet as approver.
func (r *BenefitRegistry) Approve(id uint64, actor string) error {
	if actor == "" {
		return fmt.Errorf("approver address required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	benefit, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	benefit.Approvers[actor] = struct{}{}
	if benefit.Status != BenefitStatusClaimed {
		benefit.Status = BenefitStatusApproved
	}
	benefit.UpdatedAt = time.Now().UTC()
	benefit.Audit = append(benefit.Audit, BenefitEvent{Timestamp: benefit.UpdatedAt, Actor: actor, Action: "approve"})
	return nil
}

// Claim marks the benefit as claimed without wallet verification.
func (r *BenefitRegistry) Claim(id uint64) error {
	return r.ClaimWithWallet(id, "")
}

// ClaimWithWallet marks the benefit as claimed verifying the actor.
func (r *BenefitRegistry) ClaimWithWallet(id uint64, actor string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	benefit, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	if benefit.Claimed {
		return errors.New("benefit already claimed")
	}
	if actor != "" {
		if actor != benefit.Recipient {
			if _, ok := benefit.Approvers[actor]; !ok {
				return errors.New("unauthorised claim")
			}
		}
	}
	benefit.Claimed = true
	benefit.Status = BenefitStatusClaimed
	benefit.UpdatedAt = time.Now().UTC()
	benefit.Audit = append(benefit.Audit, BenefitEvent{Timestamp: benefit.UpdatedAt, Actor: actor, Action: "claim"})
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

// StatusSummary returns counts per lifecycle bucket.
func (r *BenefitRegistry) StatusSummary() map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	counts := map[string]int{
		string(BenefitStatusPending):  0,
		string(BenefitStatusApproved): 0,
		string(BenefitStatusClaimed):  0,
		"total":                       len(r.benefits),
	}
	for _, b := range r.benefits {
		if len(b.Approvers) == 0 && !b.Claimed {
			counts[string(BenefitStatusPending)]++
		}
		if len(b.Approvers) > 0 {
			counts[string(BenefitStatusApproved)]++
		}
		if b.Claimed {
			counts[string(BenefitStatusClaimed)]++
		}
	}
	return counts
}

// AuditTrail returns the audit events for a benefit.
func (r *BenefitRegistry) AuditTrail(id uint64) ([]BenefitEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	benefit, ok := r.benefits[id]
	if !ok {
		return nil, errors.New("benefit not found")
	}
	events := make([]BenefitEvent, len(benefit.Audit))
	copy(events, benefit.Audit)
	return events, nil
}
