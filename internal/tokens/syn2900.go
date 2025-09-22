package tokens

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrInsurancePolicyNotFound = errors.New("tokens: insurance policy not found")
	ErrInsurancePolicyInactive = errors.New("tokens: insurance policy inactive")
	ErrInsuranceClaimNotFound  = errors.New("tokens: insurance claim not found")
	ErrInsuranceLimitExceeded  = errors.New("tokens: claim exceeds policy limit")
	ErrInsurancePremiumInvalid = errors.New("tokens: premium payment must be greater than zero")
)

// InsurancePolicy defines metadata for a SYN2900 general insurance token.
type InsurancePolicy struct {
	PolicyID    string
	Holder      string
	Coverage    string
	Premium     uint64
	Payout      uint64
	Deductible  uint64
	Limit       uint64
	Start       time.Time
	End         time.Time
	Claims      []ClaimRecord
	Active      bool
	PaidPremium uint64
	LastPremium time.Time
	Settled     uint64
	Reserved    uint64
}

// AvailableCoverage returns the remaining coverage capacity on the policy.
func (p *InsurancePolicy) AvailableCoverage() uint64 {
	used := p.Settled + p.Reserved
	if p.Limit <= used {
		return 0
	}
	return p.Limit - used
}

// ClaimRecord captures claim details for insurance policies.
type ClaimRecord struct {
	ClaimID string
	Amount  uint64
	Desc    string
	Time    time.Time
	Settled bool
}

// InsuranceRegistry manages SYN2900 policies.
type InsuranceRegistry struct {
	mu       sync.RWMutex
	policies map[string]*InsurancePolicy
	counter  uint64
}

// NewInsuranceRegistry creates an empty registry.
func NewInsuranceRegistry() *InsuranceRegistry {
	return &InsuranceRegistry{policies: make(map[string]*InsurancePolicy)}
}

// IssuePolicy issues a new insurance policy.
// It validates provided fields and ensures the policy period is valid.
func (r *InsuranceRegistry) IssuePolicy(holder, coverage string, premium, payout, deductible, limit uint64, start, end time.Time) (*InsurancePolicy, error) {
	if holder == "" || coverage == "" {
		return nil, errors.New("holder and coverage required")
	}
	if premium == 0 || payout == 0 {
		return nil, errors.New("premium and payout must be > 0")
	}
	if !end.After(start) || time.Now().After(end) {
		return nil, errors.New("invalid policy period")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	id := fmt.Sprintf("IP-%d", r.counter)
	p := &InsurancePolicy{
		PolicyID:    id,
		Holder:      holder,
		Coverage:    coverage,
		Premium:     premium,
		Payout:      payout,
		Deductible:  deductible,
		Limit:       limit,
		Start:       start,
		End:         end,
		Active:      true,
		LastPremium: start,
	}
	r.policies[id] = p
	return cloneInsurancePolicy(p), nil
}

// PayPremium records a premium payment for a policy.
func (r *InsuranceRegistry) PayPremium(policyID string, amount uint64) error {
	if amount == 0 {
		return ErrInsurancePremiumInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[policyID]
	if !ok {
		return ErrInsurancePolicyNotFound
	}
	if !p.Active || time.Now().After(p.End) {
		p.Active = false
		return ErrInsurancePolicyInactive
	}
	p.PaidPremium += amount
	p.LastPremium = time.Now()
	return nil
}

// FileClaim records a claim against a policy.
func (r *InsuranceRegistry) FileClaim(policyID, desc string, amount uint64) (*ClaimRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[policyID]
	if !ok {
		return nil, ErrInsurancePolicyNotFound
	}
	if !p.Active || time.Now().After(p.End) {
		p.Active = false
		return nil, ErrInsurancePolicyInactive
	}
	if amount > p.AvailableCoverage() {
		return nil, ErrInsuranceLimitExceeded
	}
	r.counter++
	id := fmt.Sprintf("IC-%d", r.counter)
	c := ClaimRecord{ClaimID: id, Amount: amount, Desc: desc, Time: time.Now()}
	p.Claims = append(p.Claims, c)
	p.Reserved += amount
	return &c, nil
}

// SettleClaim marks a claim as settled and adjusts the available coverage.
func (r *InsuranceRegistry) SettleClaim(policyID, claimID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[policyID]
	if !ok {
		return ErrInsurancePolicyNotFound
	}
	for i := range p.Claims {
		if p.Claims[i].ClaimID == claimID {
			if p.Claims[i].Settled {
				return nil
			}
			p.Claims[i].Settled = true
			if p.Reserved >= p.Claims[i].Amount {
				p.Reserved -= p.Claims[i].Amount
			} else {
				p.Reserved = 0
			}
			p.Settled += p.Claims[i].Amount
			if p.Settled >= p.Limit {
				p.Active = false
			}
			return nil
		}
	}
	return ErrInsuranceClaimNotFound
}

// GetPolicy retrieves a policy by ID.
func (r *InsuranceRegistry) GetPolicy(policyID string) (*InsurancePolicy, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.policies[policyID]
	if !ok {
		return nil, false
	}
	return cloneInsurancePolicy(p), true
}

// ListPolicies lists all insurance policies.
func (r *InsuranceRegistry) ListPolicies() []*InsurancePolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*InsurancePolicy, 0, len(r.policies))
	for _, p := range r.policies {
		res = append(res, cloneInsurancePolicy(p))
	}
	return res
}

// ListActivePolicies returns only policies that remain active.
func (r *InsuranceRegistry) ListActivePolicies(now time.Time) []*InsurancePolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*InsurancePolicy, 0)
	for _, p := range r.policies {
		if !p.Active || now.After(p.End) {
			continue
		}
		res = append(res, cloneInsurancePolicy(p))
	}
	return res
}

// TotalExposure returns the remaining exposure across active policies.
func (r *InsuranceRegistry) TotalExposure(now time.Time) uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var total uint64
	for _, p := range r.policies {
		if !p.Active || now.After(p.End) {
			continue
		}
		total += p.AvailableCoverage()
	}
	return total
}

// Deactivate marks a policy as inactive.
func (r *InsuranceRegistry) Deactivate(policyID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[policyID]
	if !ok {
		return ErrInsurancePolicyNotFound
	}
	p.Active = false
	return nil
}

func cloneInsurancePolicy(p *InsurancePolicy) *InsurancePolicy {
	cp := *p
	cp.Claims = append([]ClaimRecord(nil), p.Claims...)
	return &cp
}
