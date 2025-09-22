package tokens

import (
	"crypto/ed25519"
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
	Verifier  string
	RecordID  string
	Status    string
	Time      time.Time
	Signature []byte
}

// CarbonRegistry manages carbon credit projects and issuance.
type CarbonRegistry struct {
	mu           sync.RWMutex
	projects     map[string]*CarbonProject
	counter      uint64
	verifierKeys map[string]ed25519.PublicKey
	events       map[string][]ProjectEvent
}

// NewCarbonRegistry creates an empty carbon credit registry.
func NewCarbonRegistry() *CarbonRegistry {
	return &CarbonRegistry{
		projects:     make(map[string]*CarbonProject),
		verifierKeys: make(map[string]ed25519.PublicKey),
		events:       make(map[string][]ProjectEvent),
	}
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
	r.appendEvent(projectID, ProjectEvent{Type: "issue", Holder: holder, Amount: amount, Timestamp: time.Now()})
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
	r.appendEvent(projectID, ProjectEvent{Type: "retire", Holder: holder, Amount: amount, Timestamp: time.Now()})
	return nil
}

// RegisterVerifier associates a verification key with a verifier identity.
func (r *CarbonRegistry) RegisterVerifier(verifier string, key ed25519.PublicKey) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if key != nil {
		cp := make(ed25519.PublicKey, len(key))
		copy(cp, key)
		r.verifierKeys[verifier] = cp
	}
}

// AddSignedVerification attaches a verification record validated against the
// registered verifier key if one exists.
func (r *CarbonRegistry) AddSignedVerification(projectID, verifier, verID, status string, signature []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.projects[projectID]
	if !ok {
		return errors.New("project not found")
	}
	if key, ok := r.verifierKeys[verifier]; ok {
		payload := []byte(fmt.Sprintf("%s|%s|%s|%s", projectID, verifier, verID, status))
		if len(signature) == 0 || !ed25519.Verify(key, payload, signature) {
			return errors.New("invalid verifier signature")
		}
	}
	sigCopy := make([]byte, len(signature))
	copy(sigCopy, signature)
	v := Verification{Verifier: verifier, RecordID: verID, Status: status, Time: time.Now(), Signature: sigCopy}
	p.Verifications = append(p.Verifications, v)
	r.appendEvent(projectID, ProjectEvent{Type: "verification", Holder: verifier, Metadata: map[string]string{"record": verID, "status": status}, Timestamp: v.Time})
	return nil
}

// AddVerification retains backwards compatibility by accepting unsigned
// submissions.
func (r *CarbonRegistry) AddVerification(projectID, verifier, verID, status string) error {
	return r.AddSignedVerification(projectID, verifier, verID, status, nil)
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

// Transfer reallocates credits between holders within a project.
func (r *CarbonRegistry) Transfer(projectID, from, to string, amount uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.projects[projectID]
	if !ok {
		return errors.New("project not found")
	}
	if p.Balances[from] < amount {
		return errors.New("insufficient holder balance")
	}
	p.Balances[from] -= amount
	p.Balances[to] += amount
	r.appendEvent(projectID, ProjectEvent{Type: "transfer", Holder: from, Counterparty: to, Amount: amount, Timestamp: time.Now()})
	return nil
}

// ProjectEvent documents actions applied to a project for audit purposes.
type ProjectEvent struct {
	Type         string
	Holder       string
	Counterparty string
	Amount       uint64
	Metadata     map[string]string
	Timestamp    time.Time
}

// appendEvent appends an audit trail entry with defensive copies.
func (r *CarbonRegistry) appendEvent(projectID string, evt ProjectEvent) {
	if evt.Metadata != nil {
		cp := make(map[string]string, len(evt.Metadata))
		for k, v := range evt.Metadata {
			cp[k] = v
		}
		evt.Metadata = cp
	}
	r.events[projectID] = append(r.events[projectID], evt)
}

// ProjectTimeline exposes immutable audit trail entries.
func (r *CarbonRegistry) ProjectTimeline(projectID string, limit int) ([]ProjectEvent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	events, ok := r.events[projectID]
	if !ok {
		return nil, false
	}
	if limit <= 0 || limit >= len(events) {
		out := make([]ProjectEvent, len(events))
		copy(out, events)
		return out, true
	}
	out := make([]ProjectEvent, limit)
	copy(out, events[len(events)-limit:])
	return out, true
}

// ProjectSnapshot summarises project balances for analytics pipelines.
type ProjectSnapshot struct {
	ProjectID      string
	TotalCredits   uint64
	IssuedCredits  uint64
	RetiredCredits uint64
	HolderBalances map[string]uint64
}

// Snapshot returns a deterministic representation of the project state.
func (r *CarbonRegistry) Snapshot(projectID string) (ProjectSnapshot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.projects[projectID]
	if !ok {
		return ProjectSnapshot{}, errors.New("project not found")
	}
	balances := make(map[string]uint64, len(p.Balances))
	for holder, bal := range p.Balances {
		balances[holder] = bal
	}
	return ProjectSnapshot{
		ProjectID:      projectID,
		TotalCredits:   p.TotalCredits,
		IssuedCredits:  p.IssuedCredits,
		RetiredCredits: p.RetiredCredits,
		HolderBalances: balances,
	}, nil
}
