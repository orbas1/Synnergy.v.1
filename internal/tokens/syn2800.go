package tokens

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// LifePolicy defines the metadata for a SYN2800 life insurance token.
type LifePolicy struct {
	PolicyID    string
	Insured     string
	Beneficiary string
	Coverage    uint64
	Premium     uint64
	Start       time.Time
	End         time.Time
	PaidPremium uint64
	Claims      []Claim
	Active      bool
}

// Claim represents a filed claim against a policy.
type Claim struct {
	ClaimID string
	Amount  uint64
	Time    time.Time
	Settled bool
}

// LifePolicyRegistry manages life insurance policies.
type LifePolicyRegistry struct {
	mu       sync.RWMutex
	policies map[string]*LifePolicy
	counter  uint64
}

// NewLifePolicyRegistry creates an empty registry.
func NewLifePolicyRegistry() *LifePolicyRegistry {
	return &LifePolicyRegistry{policies: make(map[string]*LifePolicy)}
}

// IssuePolicy issues a new life insurance policy.
// It validates required fields and ensures the policy has not already expired.
func (r *LifePolicyRegistry) IssuePolicy(insured, beneficiary string, coverage, premium uint64, start, end time.Time) (*LifePolicy, error) {
	if insured == "" || beneficiary == "" {
		return nil, errors.New("insured and beneficiary required")
	}
	if coverage == 0 || premium == 0 {
		return nil, errors.New("coverage and premium must be > 0")
	}
	if !end.After(start) || time.Now().After(end) {
		return nil, errors.New("invalid policy period")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	id := fmt.Sprintf("LP-%d", r.counter)
	p := &LifePolicy{
		PolicyID:    id,
		Insured:     insured,
		Beneficiary: beneficiary,
		Coverage:    coverage,
		Premium:     premium,
		Start:       start,
		End:         end,
		Active:      true,
	}
	r.policies[id] = p
	return p, nil
}

// PayPremium records a premium payment against a policy.
func (r *LifePolicyRegistry) PayPremium(policyID string, amount uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[policyID]
	if !ok {
		return errors.New("policy not found")
	}
	if !p.Active || time.Now().After(p.End) {
		p.Active = false
		return errors.New("policy inactive")
	}
	p.PaidPremium += amount
	return nil
}

// FileClaim creates a claim record for a policy.
func (r *LifePolicyRegistry) FileClaim(policyID string, amount uint64) (*Claim, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[policyID]
	if !ok {
		return nil, errors.New("policy not found")
	}
	if !p.Active || time.Now().After(p.End) {
		p.Active = false
		return nil, errors.New("policy inactive")
	}
	r.counter++
	id := fmt.Sprintf("CL-%d", r.counter)
	c := Claim{ClaimID: id, Amount: amount, Time: time.Now()}
	p.Claims = append(p.Claims, c)
	return &c, nil
}

// GetPolicy retrieves policy information.
func (r *LifePolicyRegistry) GetPolicy(policyID string) (*LifePolicy, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.policies[policyID]
	if !ok {
		return nil, false
	}
	cp := *p
	cp.Claims = append([]Claim(nil), p.Claims...)
	return &cp, true
}

// ListPolicies lists all life policies.
func (r *LifePolicyRegistry) ListPolicies() []*LifePolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*LifePolicy, 0, len(r.policies))
	for _, p := range r.policies {
		cp := *p
		cp.Claims = append([]Claim(nil), p.Claims...)
		res = append(res, &cp)
	}
	return res
}

// Deactivate marks a policy as inactive.
func (r *LifePolicyRegistry) Deactivate(policyID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[policyID]
	if !ok {
		return errors.New("policy not found")
	}
	p.Active = false
	return nil
}
