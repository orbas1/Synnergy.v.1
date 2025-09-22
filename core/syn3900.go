package core

import (
	"errors"
	"sync"
)

// BenefitStatus enumerates lifecycle states for SYN3900 benefits.
type BenefitStatus string

const (
	// BenefitStatusPending indicates a benefit has been registered but not yet claimed.
	BenefitStatusPending BenefitStatus = "PENDING"
	// BenefitStatusClaimed marks a benefit that has been claimed by the recipient.
	BenefitStatusClaimed BenefitStatus = "CLAIMED"
)

// BenefitRecord holds metadata for a government benefit token issuance.
type BenefitRecord struct {
	ID        uint64        `json:"id"`
	Recipient string        `json:"recipient"`
	Program   string        `json:"program"`
	Amount    uint64        `json:"amount"`
	Claimed   bool          `json:"claimed"`
	Status    BenefitStatus `json:"status"`
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
	r.benefits[id] = &BenefitRecord{ID: id, Recipient: recipient, Program: program, Amount: amount, Status: BenefitStatusPending}
	return id
}

// Claim marks the benefit as claimed.
func (r *BenefitRegistry) Claim(id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	if b.Claimed {
		return errors.New("benefit already claimed")
	}
	b.Claimed = true
	b.Status = BenefitStatusClaimed
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
	cp := *b
	return &cp, true
}
