package core

import (
	"errors"
	"time"
)

// InvestmentRecord stores state for a SYN1401 investment token issuance.
type InvestmentRecord struct {
	ID          string
	Owner       string
	Principal   uint64
	Rate        float64
	Maturity    time.Time
	Accrued     uint64
	LastAccrued time.Time
}

// InvestmentRegistry manages issued investments.
type InvestmentRegistry struct {
	investments map[string]*InvestmentRecord
}

// NewInvestmentRegistry creates a new registry.
func NewInvestmentRegistry() *InvestmentRegistry {
	return &InvestmentRegistry{investments: make(map[string]*InvestmentRecord)}
}

// Issue creates a new investment record.
func (r *InvestmentRegistry) Issue(id, owner string, principal uint64, rate float64, maturity time.Time) (*InvestmentRecord, error) {
	if _, exists := r.investments[id]; exists {
		return nil, errors.New("investment already exists")
	}
	rec := &InvestmentRecord{ID: id, Owner: owner, Principal: principal, Rate: rate, Maturity: maturity, LastAccrued: time.Now()}
	r.investments[id] = rec
	return rec, nil
}

// Accrue accrues interest up to now and returns the accrued amount.
func (r *InvestmentRegistry) Accrue(id string, now time.Time) (uint64, error) {
	rec, ok := r.investments[id]
	if !ok {
		return 0, errors.New("investment not found")
	}
	if now.Before(rec.LastAccrued) {
		return 0, nil
	}
	elapsed := now.Sub(rec.LastAccrued).Hours() / 24 / 365
	interest := uint64(float64(rec.Principal) * rec.Rate * elapsed)
	rec.Accrued += interest
	rec.LastAccrued = now
	return interest, nil
}

// Redeem settles a matured investment and removes it from the registry.
// It returns the principal plus any accrued interest.
func (r *InvestmentRegistry) Redeem(id, to string, now time.Time) (uint64, error) {
	rec, ok := r.investments[id]
	if !ok {
		return 0, errors.New("investment not found")
	}
	if rec.Owner != to {
		return 0, errors.New("unauthorised redeemer")
	}
	if now.Before(rec.Maturity) {
		return 0, errors.New("investment not matured")
	}
	if _, err := r.Accrue(id, now); err != nil {
		return 0, err
	}
	total := rec.Principal + rec.Accrued
	delete(r.investments, id)
	return total, nil
}

// Get returns an investment record.
func (r *InvestmentRegistry) Get(id string) (*InvestmentRecord, bool) {
	rec, ok := r.investments[id]
	return rec, ok
}
