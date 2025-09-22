package tokens

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrPolicyNotFound is returned when a policy lookup fails.
	ErrPolicyNotFound = errors.New("tokens: policy not found")
	// ErrPolicyInactive is returned when an operation is attempted on an inactive policy.
	ErrPolicyInactive = errors.New("tokens: policy inactive")
	// ErrClaimNotFound indicates a claim lookup failure.
	ErrClaimNotFound = errors.New("tokens: claim not found")
	// ErrCoverageExceeded is raised when a claim exceeds the policy coverage.
	ErrCoverageExceeded = errors.New("tokens: claim exceeds coverage")
	// ErrInvalidPremium guards against zero payments.
	ErrInvalidPremium = errors.New("tokens: premium payment must be greater than zero")
)

// LifePolicy defines the metadata for a SYN2800 life insurance token.
type LifePolicy struct {
	PolicyID     string
	Insured      string
	Beneficiary  string
	Coverage     uint64
	Premium      uint64
	Start        time.Time
	End          time.Time
	GracePeriod  time.Duration
	PaidPremium  uint64
	LastPremium  time.Time
	Claims       []Claim
	Active       bool
	SettledClaim uint64
}

// InGracePeriod reports whether the policy is still in its grace period at the
// provided time instant.
func (p *LifePolicy) InGracePeriod(now time.Time) bool {
	if !p.Active {
		return false
	}
	if p.LastPremium.IsZero() {
		return now.Before(p.Start.Add(p.GracePeriod))
	}
	return now.Before(p.LastPremium.Add(p.GracePeriod)) && now.Before(p.End)
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
	mu           sync.RWMutex
	policies     map[string]*LifePolicy
	counter      uint64
	defaultGrace time.Duration
}

// NewLifePolicyRegistry creates an empty registry with a 30 day grace period by default.
func NewLifePolicyRegistry() *LifePolicyRegistry {
	return &LifePolicyRegistry{policies: make(map[string]*LifePolicy), defaultGrace: 30 * 24 * time.Hour}
}

// SetDefaultGrace updates the default grace period for newly issued policies.
func (r *LifePolicyRegistry) SetDefaultGrace(period time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.defaultGrace = period
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
		GracePeriod: r.defaultGrace,
		Active:      true,
		LastPremium: start,
	}
	r.policies[id] = p
	return cloneLifePolicy(p), nil
}

// PayPremium records a premium payment against a policy.
func (r *LifePolicyRegistry) PayPremium(policyID string, amount uint64) error {
	if amount == 0 {
		return ErrInvalidPremium
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[policyID]
	if !ok {
		return ErrPolicyNotFound
	}
	if !p.Active || time.Now().After(p.End) {
		p.Active = false
		return ErrPolicyInactive
	}
	p.PaidPremium += amount
	p.LastPremium = time.Now()
	return nil
}

// FileClaim creates a claim record for a policy.
func (r *LifePolicyRegistry) FileClaim(policyID string, amount uint64) (*Claim, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[policyID]
	if !ok {
		return nil, ErrPolicyNotFound
	}
	if !p.Active || time.Now().After(p.End) {
		p.Active = false
		return nil, ErrPolicyInactive
	}
	if amount > p.Coverage {
		return nil, ErrCoverageExceeded
	}
	r.counter++
	id := fmt.Sprintf("CL-%d", r.counter)
	c := Claim{ClaimID: id, Amount: amount, Time: time.Now()}
	p.Claims = append(p.Claims, c)
	return &c, nil
}

// SettleClaim marks a claim as settled and records the payout.
func (r *LifePolicyRegistry) SettleClaim(policyID, claimID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[policyID]
	if !ok {
		return ErrPolicyNotFound
	}
	for i := range p.Claims {
		if p.Claims[i].ClaimID == claimID {
			if p.Claims[i].Settled {
				return nil
			}
			p.Claims[i].Settled = true
			p.SettledClaim += p.Claims[i].Amount
			if p.SettledClaim >= p.Coverage {
				p.Active = false
			}
			return nil
		}
	}
	return ErrClaimNotFound
}

// GetPolicy retrieves policy information.
func (r *LifePolicyRegistry) GetPolicy(policyID string) (*LifePolicy, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.policies[policyID]
	if !ok {
		return nil, false
	}
	return cloneLifePolicy(p), true
}

// ListPolicies lists all life policies.
func (r *LifePolicyRegistry) ListPolicies() []*LifePolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*LifePolicy, 0, len(r.policies))
	for _, p := range r.policies {
		res = append(res, cloneLifePolicy(p))
	}
	return res
}

// ListActivePolicies returns only policies that are still active.
func (r *LifePolicyRegistry) ListActivePolicies(now time.Time) []*LifePolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*LifePolicy, 0)
	for _, p := range r.policies {
		if !p.Active || now.After(p.End) {
			continue
		}
		res = append(res, cloneLifePolicy(p))
	}
	return res
}

// Deactivate marks a policy as inactive.
func (r *LifePolicyRegistry) Deactivate(policyID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[policyID]
	if !ok {
		return ErrPolicyNotFound
	}
	p.Active = false
	return nil
}

func cloneLifePolicy(p *LifePolicy) *LifePolicy {
	cp := *p
	cp.Claims = append([]Claim(nil), p.Claims...)
	return &cp
}
