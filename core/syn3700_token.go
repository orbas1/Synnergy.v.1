package core

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// Component represents a token and its weight within the index.
type Component struct {
	Token  string
	Weight float64
}

type indexComponent struct {
	Component
	Drift         float64
	CurrentWeight float64
	UpdatedAt     time.Time
}

// SYN3700AuditEntry records privileged operations performed on the index.
type SYN3700AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Token     string    `json:"token,omitempty"`
	Previous  float64   `json:"previous_weight,omitempty"`
	Weight    float64   `json:"weight,omitempty"`
	Drift     float64   `json:"drift,omitempty"`
}

// SYN3700Status summarises telemetry for dashboards and CLI.
type SYN3700Status struct {
	ComponentCount  int       `json:"component_count"`
	ControllerCount int       `json:"controller_count"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SYN3700Snapshot exposes a serialisable view for JSON output.
type SYN3700Snapshot struct {
	Name        string      `json:"name"`
	Symbol      string      `json:"symbol"`
	Components  []Component `json:"components"`
	Controllers []string    `json:"controllers"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// SYN3700Token models a weighted index token with controller governance.
type SYN3700Token struct {
	mu          sync.RWMutex
	Name        string
	Symbol      string
	components  map[string]*indexComponent
	controllers map[string]struct{}
	audit       []SYN3700AuditEntry
	updatedAt   time.Time
}

// NewSYN3700Token creates an empty index token.
func NewSYN3700Token(name, symbol string) *SYN3700Token {
	return &SYN3700Token{
		Name:        name,
		Symbol:      symbol,
		components:  make(map[string]*indexComponent),
		controllers: make(map[string]struct{}),
		updatedAt:   time.Now().UTC(),
	}
}

// RegisterController grants privileged access to an operator address.
func (t *SYN3700Token) RegisterController(address string) {
	if address == "" {
		return
	}
	t.mu.Lock()
	t.controllers[address] = struct{}{}
	t.updatedAt = time.Now().UTC()
	t.audit = append(t.audit, SYN3700AuditEntry{
		Timestamp: t.updatedAt,
		Actor:     address,
		Action:    "register_controller",
	})
	t.mu.Unlock()
}

// Controllers returns the registered controller addresses sorted for determinism.
func (t *SYN3700Token) Controllers() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	res := make([]string, 0, len(t.controllers))
	for addr := range t.controllers {
		res = append(res, addr)
	}
	sort.Strings(res)
	return res
}

func (t *SYN3700Token) isController(address string) bool {
	if address == "" {
		return false
	}
	_, ok := t.controllers[address]
	return ok
}

func (t *SYN3700Token) requireController(actor string) error {
	if len(t.controllers) == 0 {
		return nil
	}
	if !t.isController(actor) {
		return fmt.Errorf("controller required: %s", actor)
	}
	return nil
}

// AddComponent adds or updates a component token with a given weight.
func (t *SYN3700Token) AddComponent(token string, weight float64) {
	_ = t.AddComponentControlled("", token, weight, 0)
}

// AddComponentControlled enforces controller permissions and drift policies.
func (t *SYN3700Token) AddComponentControlled(actor, token string, weight, drift float64) error {
	if weight <= 0 {
		return errors.New("invalid weight")
	}
	if drift < 0 || drift > 1 {
		return errors.New("component drift out of range")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if err := t.requireController(actor); err != nil {
		return err
	}
	if token == "" {
		return errors.New("token symbol required")
	}
	now := time.Now().UTC()
	entry := &indexComponent{
		Component: Component{Token: token, Weight: weight},
		Drift:     drift,
		UpdatedAt: now,
	}
	entry.CurrentWeight = weight
	prev := t.components[token]
	t.components[token] = entry
	t.updatedAt = now
	prevWeight := 0.0
	if prev != nil {
		prevWeight = prev.Weight
	}
	t.audit = append(t.audit, SYN3700AuditEntry{
		Timestamp: now,
		Actor:     actor,
		Action:    "add_component",
		Token:     token,
		Previous:  prevWeight,
		Weight:    weight,
		Drift:     drift,
	})
	return nil
}

// RemoveComponent deletes a component from the index without controller enforcement.
func (t *SYN3700Token) RemoveComponent(token string) error {
	return t.RemoveComponentControlled("", token)
}

// RemoveComponentControlled removes a component after controller validation.
func (t *SYN3700Token) RemoveComponentControlled(actor, token string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if err := t.requireController(actor); err != nil {
		return err
	}
	comp, ok := t.components[token]
	if !ok {
		return errors.New("component not found")
	}
	delete(t.components, token)
	now := time.Now().UTC()
	t.updatedAt = now
	t.audit = append(t.audit, SYN3700AuditEntry{
		Timestamp: now,
		Actor:     actor,
		Action:    "remove_component",
		Token:     token,
		Previous:  comp.Weight,
	})
	return nil
}

// ListComponents returns a snapshot of all components.
func (t *SYN3700Token) ListComponents() []Component {
	t.mu.RLock()
	defer t.mu.RUnlock()
	comps := make([]Component, 0, len(t.components))
	for _, comp := range t.components {
		comps = append(comps, Component{Token: comp.Token, Weight: comp.Weight})
	}
	sort.Slice(comps, func(i, j int) bool { return comps[i].Token < comps[j].Token })
	return comps
}

// Snapshot returns a serialisable state summary for CLI/API usage.
func (t *SYN3700Token) Snapshot() SYN3700Snapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()
	comps := make([]Component, 0, len(t.components))
	for _, comp := range t.components {
		comps = append(comps, Component{Token: comp.Token, Weight: comp.Weight})
	}
	sort.Slice(comps, func(i, j int) bool { return comps[i].Token < comps[j].Token })
	controllers := make([]string, 0, len(t.controllers))
	for addr := range t.controllers {
		controllers = append(controllers, addr)
	}
	sort.Strings(controllers)
	return SYN3700Snapshot{
		Name:        t.Name,
		Symbol:      t.Symbol,
		Components:  comps,
		Controllers: controllers,
		UpdatedAt:   t.updatedAt,
	}
}

// Status returns operational telemetry for dashboards and CLI consumers.
func (t *SYN3700Token) Status() SYN3700Status {
	snap := t.Snapshot()
	return SYN3700Status{
		ComponentCount:  len(snap.Components),
		ControllerCount: len(snap.Controllers),
		UpdatedAt:       snap.UpdatedAt,
	}
}

// RebalanceControlled snaps all components back to their target weights.
func (t *SYN3700Token) RebalanceControlled(actor string) (map[string][2]float64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if err := t.requireController(actor); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	changes := make(map[string][2]float64, len(t.components))
	for token, comp := range t.components {
		prev := comp.CurrentWeight
		comp.CurrentWeight = comp.Weight
		comp.UpdatedAt = now
		changes[token] = [2]float64{comp.Weight, prev}
		t.audit = append(t.audit, SYN3700AuditEntry{
			Timestamp: now,
			Actor:     actor,
			Action:    "rebalance",
			Token:     token,
			Previous:  prev,
			Weight:    comp.Weight,
			Drift:     comp.Drift,
		})
	}
	t.updatedAt = now
	return changes, nil
}

// AuditTrail returns a copy of recorded administrative actions.
func (t *SYN3700Token) AuditTrail() []SYN3700AuditEntry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]SYN3700AuditEntry, len(t.audit))
	copy(out, t.audit)
	return out
}

// Value computes the index value given a map of token prices.
func (t *SYN3700Token) Value(prices map[string]float64) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var total float64
	for token, comp := range t.components {
		total += comp.CurrentWeight * prices[token]
	}
	return total
}
