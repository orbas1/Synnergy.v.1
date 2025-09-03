package tokens

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// InsurancePolicy defines metadata for a SYN2900 general insurance token.
type InsurancePolicy struct {
	PolicyID   string
	Holder     string
	Coverage   string
	Premium    uint64
	Payout     uint64
	Deductible uint64
	Limit      uint64
	Start      time.Time
	End        time.Time
	Claims     []ClaimRecord
	Active     bool
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
		PolicyID:   id,
		Holder:     holder,
		Coverage:   coverage,
		Premium:    premium,
		Payout:     payout,
		Deductible: deductible,
		Limit:      limit,
		Start:      start,
		End:        end,
		Active:     true,
	}
	r.policies[id] = p
	return p, nil
}

// FileClaim records a claim against a policy.
func (r *InsuranceRegistry) FileClaim(policyID, desc string, amount uint64) (*ClaimRecord, error) {
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
	id := fmt.Sprintf("IC-%d", r.counter)
	c := ClaimRecord{ClaimID: id, Amount: amount, Desc: desc, Time: time.Now()}
	p.Claims = append(p.Claims, c)
	return &c, nil
}

// GetPolicy retrieves a policy by ID.
func (r *InsuranceRegistry) GetPolicy(policyID string) (*InsurancePolicy, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.policies[policyID]
	if !ok {
		return nil, false
	}
	cp := *p
	cp.Claims = append([]ClaimRecord(nil), p.Claims...)
	return &cp, true
}

// ListPolicies lists all insurance policies.
func (r *InsuranceRegistry) ListPolicies() []*InsurancePolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*InsurancePolicy, 0, len(r.policies))
	for _, p := range r.policies {
		cp := *p
		cp.Claims = append([]ClaimRecord(nil), p.Claims...)
		res = append(res, &cp)
	}
	return res
}

// Deactivate marks a policy as inactive.
func (r *InsuranceRegistry) Deactivate(policyID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[policyID]
	if !ok {
		return errors.New("policy not found")
	}
	p.Active = false
	return nil
}
