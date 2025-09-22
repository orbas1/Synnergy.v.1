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
	Token  string  `json:"token"`
	Weight float64 `json:"weight"`
	Drift  float64 `json:"drift"`
	Target float64 `json:"target"`
}

type componentState struct {
	Component
	lastRebalance time.Time
}

// IndexAuditEntry captures actions performed on the SYN3700 token.
type IndexAuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Token     string    `json:"token,omitempty"`
	Actor     string    `json:"actor"`
	Details   string    `json:"details,omitempty"`
}

// SYN3700Token models a weighted index token.
type SYN3700Token struct {
	mu          sync.RWMutex
	Name        string
	Symbol      string
	components  map[string]*componentState
	controllers map[string]struct{}
	audit       []IndexAuditEntry
}

// NewSYN3700Token creates an empty index token.
func NewSYN3700Token(name, symbol string) *SYN3700Token {
	return &SYN3700Token{
		Name:        name,
		Symbol:      symbol,
		components:  make(map[string]*componentState),
		controllers: make(map[string]struct{}),
	}
}

// AddController registers an address that can mutate the index.
func (t *SYN3700Token) AddController(address string) {
	if address == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.controllers[address]; exists {
		return
	}
	t.controllers[address] = struct{}{}
	t.audit = append(t.audit, IndexAuditEntry{
		Timestamp: time.Now().UTC(),
		Action:    "controller_added",
		Actor:     address,
	})
}

// HasController reports whether the address may modify the index.
func (t *SYN3700Token) HasController(address string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.controllers[address]
	return ok
}

// Controllers returns the registered controllers in deterministic order.
func (t *SYN3700Token) Controllers() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]string, 0, len(t.controllers))
	for addr := range t.controllers {
		out = append(out, addr)
	}
	sort.Strings(out)
	return out
}

// AddComponent adds or updates a component token with a given weight and drift threshold.
func (t *SYN3700Token) AddComponent(actor, token string, weight, drift float64) error {
	if token == "" {
		return errors.New("token symbol required")
	}
	if weight <= 0 {
		return errors.New("invalid weight")
	}
	if drift < 0 || drift > 1 {
		return errors.New("component drift must be between 0 and 1")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	st, ok := t.components[token]
	if !ok {
		st = &componentState{Component: Component{Token: token}}
		t.components[token] = st
	}
	st.Weight = weight
	st.Target = weight
	st.Drift = drift
	st.lastRebalance = time.Now().UTC()
	t.audit = append(t.audit, IndexAuditEntry{
		Timestamp: time.Now().UTC(),
		Action:    "component_added",
		Token:     token,
		Actor:     actor,
		Details:   fmt.Sprintf("weight=%.4f drift=%.4f", weight, drift),
	})
	return nil
}

// RemoveComponent deletes a component from the index.
func (t *SYN3700Token) RemoveComponent(actor, token string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.components[token]; !ok {
		return errors.New("component not found")
	}
	delete(t.components, token)
	t.audit = append(t.audit, IndexAuditEntry{
		Timestamp: time.Now().UTC(),
		Action:    "component_removed",
		Token:     token,
		Actor:     actor,
	})
	return nil
}

// ListComponents returns a snapshot of all components.
func (t *SYN3700Token) ListComponents() []Component {
	t.mu.RLock()
	defer t.mu.RUnlock()
	comps := make([]Component, 0, len(t.components))
	for _, st := range t.components {
		comps = append(comps, st.Component)
	}
	sort.Slice(comps, func(i, j int) bool { return comps[i].Token < comps[j].Token })
	return comps
}

// Value computes the index value given a map of token prices.
func (t *SYN3700Token) Value(prices map[string]float64) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var total float64
	for token, st := range t.components {
		price := prices[token]
		total += st.Weight * price
	}
	return total
}

// Snapshot returns a serialisable view of the index state.
type IndexSnapshot struct {
	Name       string      `json:"name"`
	Symbol     string      `json:"symbol"`
	Components []Component `json:"components"`
}

// Snapshot captures the current index state.
func (t *SYN3700Token) Snapshot() IndexSnapshot {
	comps := t.ListComponents()
	return IndexSnapshot{
		Name:       t.Name,
		Symbol:     t.Symbol,
		Components: comps,
	}
}

// Telemetry summarises controller and component counts.
type IndexTelemetry struct {
	ComponentCount  int `json:"component_count"`
	ControllerCount int `json:"controller_count"`
}

// Telemetry returns an aggregated view for dashboards.
func (t *SYN3700Token) Telemetry() IndexTelemetry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return IndexTelemetry{
		ComponentCount:  len(t.components),
		ControllerCount: len(t.controllers),
	}
}

// Rebalance resets component weights to their target and reports the change.
func (t *SYN3700Token) Rebalance(actor string) map[string][2]float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	updates := make(map[string][2]float64)
	for token, st := range t.components {
		from := st.Weight
		to := st.Target
		st.Weight = to
		st.lastRebalance = time.Now().UTC()
		updates[token] = [2]float64{from, to}
	}
	t.audit = append(t.audit, IndexAuditEntry{
		Timestamp: time.Now().UTC(),
		Action:    "rebalance",
		Actor:     actor,
		Details:   fmt.Sprintf("components=%d", len(updates)),
	})
	return updates
}

// AuditTrail returns a copy of audit entries ordered oldest to newest.
func (t *SYN3700Token) AuditTrail() []IndexAuditEntry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	audit := append([]IndexAuditEntry(nil), t.audit...)
	sort.SliceStable(audit, func(i, j int) bool { return audit[i].Timestamp.Before(audit[j].Timestamp) })
	return audit
}
