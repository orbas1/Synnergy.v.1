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
	Token         string    `json:"token"`
	TargetWeight  float64   `json:"target_weight"`
	Drift         float64   `json:"drift"`
	CurrentWeight float64   `json:"current_weight"`
	LastUpdated   time.Time `json:"last_updated"`
}

// IndexEvent captures auditable actions on the SYN3700 token.
type IndexEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Token     string    `json:"token,omitempty"`
	Details   string    `json:"details,omitempty"`
}

// SYN3700Token models a weighted index token with controller governance.
type SYN3700Token struct {
	mu          sync.RWMutex
	Name        string
	Symbol      string
	components  map[string]*Component
	controllers map[string]struct{}
	audit       []IndexEvent
}

// NewSYN3700Token creates an empty index token.
func NewSYN3700Token(name, symbol string) *SYN3700Token {
	return &SYN3700Token{
		Name:        name,
		Symbol:      symbol,
		components:  make(map[string]*Component),
		controllers: make(map[string]struct{}),
	}
}

// AddController registers an address as an index controller.
func (t *SYN3700Token) AddController(addr string) {
	if addr == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.controllers[addr]; exists {
		return
	}
	t.controllers[addr] = struct{}{}
	t.audit = append(t.audit, IndexEvent{Timestamp: time.Now().UTC(), Action: "controller_added", Details: addr})
}

func (t *SYN3700Token) isController(addr string) bool {
	if addr == "" {
		return false
	}
	_, ok := t.controllers[addr]
	return ok
}

// AddComponent adds or updates a component token with a given weight and drift tolerance.
func (t *SYN3700Token) AddComponent(token string, weight, drift float64, signer string) error {
	if token == "" {
		return errors.New("token symbol required")
	}
	if weight <= 0 {
		return fmt.Errorf("invalid weight")
	}
	if drift < 0 || drift > 1 {
		return fmt.Errorf("component drift out of range")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.controllers) > 0 {
		if _, ok := t.controllers[signer]; !ok {
			return errors.New("controller authorization required")
		}
	}
	comp := &Component{
		Token:         token,
		TargetWeight:  weight,
		Drift:         drift,
		CurrentWeight: weight,
		LastUpdated:   time.Now().UTC(),
	}
	t.components[token] = comp
	t.audit = append(t.audit, IndexEvent{Timestamp: comp.LastUpdated, Action: "component_added", Token: token, Details: fmt.Sprintf("weight=%f", weight)})
	return nil
}

// RemoveComponent deletes a component from the index.
func (t *SYN3700Token) RemoveComponent(token, signer string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.controllers) > 0 {
		if _, ok := t.controllers[signer]; !ok {
			return errors.New("controller authorization required")
		}
	}
	if _, ok := t.components[token]; !ok {
		return errors.New("component not found")
	}
	delete(t.components, token)
	t.audit = append(t.audit, IndexEvent{Timestamp: time.Now().UTC(), Action: "component_removed", Token: token})
	return nil
}

// ListComponents returns a snapshot of all components.
func (t *SYN3700Token) ListComponents() []Component {
	t.mu.RLock()
	defer t.mu.RUnlock()
	comps := make([]Component, 0, len(t.components))
	for _, c := range t.components {
		comps = append(comps, *c)
	}
	sort.Slice(comps, func(i, j int) bool { return comps[i].Token < comps[j].Token })
	return comps
}

// Value computes the index value given a map of token prices.
func (t *SYN3700Token) Value(prices map[string]float64) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var total float64
	for tok, c := range t.components {
		total += c.TargetWeight * prices[tok]
	}
	return total
}

// Controllers returns the list of registered controllers.
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

// ComponentCount returns the number of configured components.
func (t *SYN3700Token) ComponentCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.components)
}

// ControllerCount returns the number of controllers.
func (t *SYN3700Token) ControllerCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.controllers)
}

// Rebalance adjusts current weights back to target values and returns the changes applied.
func (t *SYN3700Token) Rebalance(signer string) (map[string][2]float64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.controllers) > 0 {
		if _, ok := t.controllers[signer]; !ok {
			return nil, errors.New("controller authorization required")
		}
	}
	res := make(map[string][2]float64, len(t.components))
	for token, comp := range t.components {
		before := comp.CurrentWeight
		comp.CurrentWeight = comp.TargetWeight
		comp.LastUpdated = time.Now().UTC()
		res[token] = [2]float64{before, comp.CurrentWeight}
	}
	t.audit = append(t.audit, IndexEvent{Timestamp: time.Now().UTC(), Action: "rebalance", Details: fmt.Sprintf("components=%d", len(res))})
	return res, nil
}

// Audit returns the audit trail for index operations.
func (t *SYN3700Token) Audit() []IndexEvent {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]IndexEvent, len(t.audit))
	copy(out, t.audit)
	return out
}
