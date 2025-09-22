package core

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// Component represents a token and its weight within the index.
type Component struct {
	Token       string    `json:"token"`
	Weight      float64   `json:"weight"`
	Drift       float64   `json:"drift"`
	LastUpdated time.Time `json:"last_updated"`
}

// IndexEvent captures an audit trail entry for index operations.
type IndexEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
}

// IndexTelemetry reports high-level counts for monitoring.
type IndexTelemetry struct {
	ComponentCount  int `json:"component_count"`
	ControllerCount int `json:"controller_count"`
}

// SYN3700Snapshot persists the state of the index token.
type SYN3700Snapshot struct {
	Name        string       `json:"name"`
	Symbol      string       `json:"symbol"`
	Components  []Component  `json:"components"`
	Controllers []string     `json:"controllers"`
	Audit       []IndexEvent `json:"audit"`
}

// SYN3700Token models a weighted index token with controller governance.
type SYN3700Token struct {
	mu          sync.RWMutex
	Name        string
	Symbol      string
	Components  map[string]Component
	Controllers map[string]bool
	AuditLog    []IndexEvent
}

// NewSYN3700Token creates an empty index token.
func NewSYN3700Token(name, symbol string) *SYN3700Token {
	return &SYN3700Token{
		Name:        name,
		Symbol:      symbol,
		Components:  make(map[string]Component),
		Controllers: make(map[string]bool),
	}
}

// NewSYN3700FromSnapshot restores an index token from persisted state.
func NewSYN3700FromSnapshot(snap *SYN3700Snapshot) *SYN3700Token {
	if snap == nil {
		return nil
	}
	token := NewSYN3700Token(snap.Name, snap.Symbol)
	for _, c := range snap.Components {
		token.Components[c.Token] = c
	}
	for _, ctrl := range snap.Controllers {
		token.Controllers[ctrl] = true
	}
	token.AuditLog = append([]IndexEvent{}, snap.Audit...)
	return token
}

// AddController authorises an address to mutate the index.
func (t *SYN3700Token) AddController(addr string) {
	if addr == "" {
		return
	}
	t.mu.Lock()
	if t.Controllers == nil {
		t.Controllers = make(map[string]bool)
	}
	if !t.Controllers[addr] {
		t.Controllers[addr] = true
		t.AuditLog = append(t.AuditLog, IndexEvent{Timestamp: time.Now().UTC(), Actor: addr, Action: "controller:add"})
	}
	t.mu.Unlock()
}

func (t *SYN3700Token) ensureController(addr string) error {
	if len(t.Controllers) == 0 {
		return nil
	}
	if !t.Controllers[addr] {
		return errors.New("controller required")
	}
	return nil
}

// AddComponent adds or updates a component token with a given weight and drift tolerance.
func (t *SYN3700Token) AddComponent(token string, weight, drift float64, actor string) error {
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
	if err := t.ensureController(actor); err != nil {
		return err
	}
	if t.Components == nil {
		t.Components = make(map[string]Component)
	}
	comp := Component{Token: token, Weight: weight, Drift: drift, LastUpdated: time.Now().UTC()}
	t.Components[token] = comp
	t.AuditLog = append(t.AuditLog, IndexEvent{Timestamp: comp.LastUpdated, Actor: actor, Action: "component:add", Details: token})
	return nil
}

// RemoveComponent deletes a component from the index.
func (t *SYN3700Token) RemoveComponent(token string, actor string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if err := t.ensureController(actor); err != nil {
		return err
	}
	if _, ok := t.Components[token]; !ok {
		return errors.New("component not found")
	}
	delete(t.Components, token)
	t.AuditLog = append(t.AuditLog, IndexEvent{Timestamp: time.Now().UTC(), Actor: actor, Action: "component:remove", Details: token})
	return nil
}

// ListComponents returns a snapshot of all components.
func (t *SYN3700Token) ListComponents() []Component {
	t.mu.RLock()
	defer t.mu.RUnlock()
	comps := make([]Component, 0, len(t.Components))
	for _, c := range t.Components {
		comps = append(comps, c)
	}
	sort.Slice(comps, func(i, j int) bool { return comps[i].Token < comps[j].Token })
	return comps
}

// Value computes the index value given a map of token prices.
func (t *SYN3700Token) Value(prices map[string]float64) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var total float64
	for token, comp := range t.Components {
		total += comp.Weight * prices[token]
	}
	return total
}

// Telemetry summarises controller/component counts.
func (t *SYN3700Token) Telemetry() IndexTelemetry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	tele := IndexTelemetry{ComponentCount: len(t.Components), ControllerCount: len(t.Controllers)}
	return tele
}

// Controllers returns the list of authorised controllers.
func (t *SYN3700Token) ControllersList() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]string, 0, len(t.Controllers))
	for addr := range t.Controllers {
		out = append(out, addr)
	}
	sort.Strings(out)
	return out
}

// Rebalance returns a report describing how weights should adjust back to targets.
func (t *SYN3700Token) Rebalance(actor string) map[string][2]float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	if err := t.ensureController(actor); err != nil {
		return map[string][2]float64{}
	}
	report := make(map[string][2]float64, len(t.Components))
	now := time.Now().UTC()
	for token, comp := range t.Components {
		current := comp.Weight
		targetHigh := comp.Weight * (1 + comp.Drift)
		report[token] = [2]float64{current, targetHigh}
		comp.LastUpdated = now
		t.Components[token] = comp
	}
	t.AuditLog = append(t.AuditLog, IndexEvent{Timestamp: now, Actor: actor, Action: "rebalance"})
	return report
}

// Audit returns the audit log entries.
func (t *SYN3700Token) Audit() []IndexEvent {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return append([]IndexEvent{}, t.AuditLog...)
}

// Snapshot returns a serialisable copy of the token state.
func (t *SYN3700Token) Snapshot() *SYN3700Snapshot {
	return &SYN3700Snapshot{
		Name:        t.Name,
		Symbol:      t.Symbol,
		Components:  t.ListComponents(),
		Controllers: t.ControllersList(),
		Audit:       t.Audit(),
	}
}
