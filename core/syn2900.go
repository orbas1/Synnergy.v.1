package core

import (
	"errors"
	"time"
)

// TokenInsurancePolicy represents a blockchain based insurance policy.
type TokenInsurancePolicy struct {
	PolicyID   string
	Holder     string
	Coverage   string
	Premium    uint64
	Payout     uint64
	Deductible uint64
	Limit      uint64
	Start      time.Time
	End        time.Time
	Claimed    bool
}

// NewTokenInsurancePolicy issues a new policy.
func NewTokenInsurancePolicy(id, holder, coverage string, premium, payout, deductible, limit uint64, start, end time.Time) *TokenInsurancePolicy {
	return &TokenInsurancePolicy{
		PolicyID:   id,
		Holder:     holder,
		Coverage:   coverage,
		Premium:    premium,
		Payout:     payout,
		Deductible: deductible,
		Limit:      limit,
		Start:      start,
		End:        end,
	}
}

// IsActive reports whether the policy is active at the given time.
func (p *TokenInsurancePolicy) IsActive(now time.Time) bool {
	return !now.Before(p.Start) && now.Before(p.End) && !p.Claimed
}

// Claim marks the policy as claimed and returns the payout.
func (p *TokenInsurancePolicy) Claim(now time.Time) (uint64, error) {
	if !p.IsActive(now) {
		return 0, errors.New("policy inactive")
	}
	p.Claimed = true
	return p.Payout, nil
}
