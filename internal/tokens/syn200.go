package tokens

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// CarbonProject represents a registered carbon offset project.
type CarbonProject struct {
	ID             string
	Owner          string
	Name           string
	TotalCredits   uint64
	IssuedCredits  uint64
	RetiredCredits uint64
	Balances       map[string]uint64
	Verifications  []Verification
	CreatedAt      time.Time
}

// Verification captures an audit or verification record for a project.
type Verification struct {
	Verifier string
	RecordID string
	Status   string
	Time     time.Time
}

// CarbonRegistry manages carbon credit projects and issuance.
type CarbonRegistry struct {
	mu       sync.RWMutex
	projects map[string]*CarbonProject
	counter  uint64
}

// NewCarbonRegistry creates an empty carbon credit registry.
func NewCarbonRegistry() *CarbonRegistry {
	return &CarbonRegistry{projects: make(map[string]*CarbonProject)}
}

// Register creates a new carbon project and returns its identifier.
func (r *CarbonRegistry) Register(owner, name string, total uint64) *CarbonProject {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	id := fmt.Sprintf("CP-%d", r.counter)
	p := &CarbonProject{
		ID:           id,
		Owner:        owner,
		Name:         name,
		TotalCredits: total,
		Balances:     make(map[string]uint64),
		CreatedAt:    time.Now(),
	}
	r.projects[id] = p
	return p
}

// Issue credits from a project to a holder.
func (r *CarbonRegistry) Issue(projectID, holder string, amount uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.projects[projectID]
	if !ok {
		return errors.New("project not found")
	}
	if p.IssuedCredits+amount > p.TotalCredits {
		return errors.New("insufficient credits remaining")
	}
	p.IssuedCredits += amount
	p.Balances[holder] += amount
	return nil
}

// Retire removes credits from circulation for a holder.
func (r *CarbonRegistry) Retire(projectID, holder string, amount uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.projects[projectID]
	if !ok {
		return errors.New("project not found")
	}
	bal := p.Balances[holder]
	if bal < amount {
		return errors.New("insufficient holder balance")
	}
	p.Balances[holder] = bal - amount
	p.RetiredCredits += amount
	p.IssuedCredits -= amount
	return nil
}

// AddVerification attaches a verification record to a project.
func (r *CarbonRegistry) AddVerification(projectID, verifier, verID, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.projects[projectID]
	if !ok {
		return errors.New("project not found")
	}
	v := Verification{Verifier: verifier, RecordID: verID, Status: status, Time: time.Now()}
	p.Verifications = append(p.Verifications, v)
	return nil
}

// Verifications lists verification records for a project.
func (r *CarbonRegistry) Verifications(projectID string) ([]Verification, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.projects[projectID]
	if !ok {
		return nil, false
	}
	return append([]Verification(nil), p.Verifications...), true
}

// ProjectInfo returns project details by ID.
func (r *CarbonRegistry) ProjectInfo(projectID string) (*CarbonProject, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.projects[projectID]
	if !ok {
		return nil, false
	}
	cp := *p
	cp.Balances = make(map[string]uint64)
	for k, v := range p.Balances {
		cp.Balances[k] = v
	}
	cp.Verifications = append([]Verification(nil), p.Verifications...)
	return &cp, true
}

// ListProjects returns all registered carbon projects.
func (r *CarbonRegistry) ListProjects() []*CarbonProject {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*CarbonProject, 0, len(r.projects))
	for _, p := range r.projects {
		cp := *p
		cp.Balances = make(map[string]uint64)
		for k, v := range p.Balances {
			cp.Balances[k] = v
		}
		cp.Verifications = append([]Verification(nil), p.Verifications...)
		res = append(res, &cp)
	}
	return res
}
