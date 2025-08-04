package core

import "errors"

// BenefitRecord holds metadata for a government benefit token issuance.
type BenefitRecord struct {
	ID        uint64
	Recipient string
	Program   string
	Amount    uint64
	Claimed   bool
}

// BenefitRegistry manages benefit records.
type BenefitRegistry struct {
	benefits map[uint64]*BenefitRecord
	nextID   uint64
}

// NewBenefitRegistry creates a new registry.
func NewBenefitRegistry() *BenefitRegistry {
	return &BenefitRegistry{benefits: make(map[uint64]*BenefitRecord)}
}

// RegisterBenefit records a new benefit and returns its ID.
func (r *BenefitRegistry) RegisterBenefit(recipient, program string, amount uint64) uint64 {
	r.nextID++
	id := r.nextID
	r.benefits[id] = &BenefitRecord{ID: id, Recipient: recipient, Program: program, Amount: amount}
	return id
}

// Claim marks the benefit as claimed.
func (r *BenefitRegistry) Claim(id uint64) error {
	b, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	if b.Claimed {
		return errors.New("benefit already claimed")
	}
	b.Claimed = true
	return nil
}

// GetBenefit retrieves a benefit by ID.
func (r *BenefitRegistry) GetBenefit(id uint64) (*BenefitRecord, bool) {
	b, ok := r.benefits[id]
	return b, ok
}
