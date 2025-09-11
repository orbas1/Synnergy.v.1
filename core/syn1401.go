package core

import (
	"errors"
	"sync"
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
	mu          sync.RWMutex
	investments map[string]*InvestmentRecord
}

// NewInvestmentRegistry creates a new registry.
func NewInvestmentRegistry() *InvestmentRegistry {
	return &InvestmentRegistry{investments: make(map[string]*InvestmentRecord)}
}

var (
	// ErrInvestmentExists indicates an ID collision when issuing a record.
	ErrInvestmentExists = errors.New("investment already exists")
	// ErrInvestmentNotFound is returned when a lookup fails.
	ErrInvestmentNotFound = errors.New("investment not found")
	// ErrUnauthorizedRedeemer signals a redeem attempt from a non-owner.
	ErrUnauthorizedRedeemer = errors.New("unauthorised redeemer")
	// ErrInvestmentNotMatured signals redeem prior to maturity.
	ErrInvestmentNotMatured = errors.New("investment not matured")
)

// Issue creates a new investment record.
func (r *InvestmentRegistry) Issue(id, owner string, principal uint64, rate float64, maturity time.Time) (*InvestmentRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.investments[id]; exists {
		return nil, ErrInvestmentExists
	}
	rec := &InvestmentRecord{ID: id, Owner: owner, Principal: principal, Rate: rate, Maturity: maturity, LastAccrued: time.Now()}
	r.investments[id] = rec
	return rec, nil
}

// Accrue accrues interest up to now and returns the accrued amount.
func (r *InvestmentRegistry) Accrue(id string, now time.Time) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.investments[id]
	if !ok {
		return 0, ErrInvestmentNotFound
	}
	return accrueInvestment(rec, now), nil
}

func accrueInvestment(rec *InvestmentRecord, now time.Time) uint64 {
	if now.Before(rec.LastAccrued) {
		return 0
	}
	elapsed := now.Sub(rec.LastAccrued).Hours() / 24 / 365
	interest := uint64(float64(rec.Principal) * rec.Rate * elapsed)
	rec.Accrued += interest
	rec.LastAccrued = now
	return interest
}

// Redeem settles a matured investment and removes it from the registry.
// It returns the principal plus any accrued interest.
func (r *InvestmentRegistry) Redeem(id, to string, now time.Time) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.investments[id]
	if !ok {
		return 0, ErrInvestmentNotFound
	}
	if rec.Owner != to {
		return 0, ErrUnauthorizedRedeemer
	}
	if now.Before(rec.Maturity) {
		return 0, ErrInvestmentNotMatured
	}
	accrueInvestment(rec, now)
	total := rec.Principal + rec.Accrued
	delete(r.investments, id)
	return total, nil
}

// Get returns an investment record.
func (r *InvestmentRegistry) Get(id string) (*InvestmentRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec, ok := r.investments[id]
	return rec, ok
}
