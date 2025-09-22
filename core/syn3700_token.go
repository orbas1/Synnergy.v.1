package core

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

var (
	// ErrControllerRequired is returned when a mutating operation is attempted
	// by an account that is not authorised to govern the index token.
	ErrControllerRequired = errors.New("controller required")
	// ErrComponentExists signals that a component already exists in the index.
	ErrComponentExists = errors.New("component already exists")
	// ErrComponentNotFound is returned when referencing an unknown component.
	ErrComponentNotFound = errors.New("component not found")
)

// Component represents a token and its current portfolio characteristics. It is
// exported so CLI and documentation tooling can serialise token state directly.
//
// Drift captures the permitted delta between the target and the current weight,
// while Target records the intended allocation after the next rebalance cycle.
// LastUpdated indicates when controllers last adjusted the component.
//
// json struct tags are included for Stage 73 persistence snapshots.
type Component struct {
	Token       string    `json:"token"`
	Weight      float64   `json:"weight"`
	Drift       float64   `json:"drift"`
	Target      float64   `json:"target"`
	LastUpdated time.Time `json:"last_updated"`
}

type syn3700Component struct {
	Component
	AddedBy        Address
	LastRebalanced time.Time
}

type syn3700Controller struct {
	Address   Address    `json:"address"`
	GrantedBy Address    `json:"granted_by"`
	GrantedAt time.Time  `json:"granted_at"`
	Revoked   *time.Time `json:"revoked_at,omitempty"`
}

// SYN3700AuditRecord captures governance events for transparency.
type SYN3700AuditRecord struct {
	Timestamp time.Time `json:"timestamp"`
	Actor     Address   `json:"actor"`
	Action    string    `json:"action"`
	Details   string    `json:"details,omitempty"`
}

// SYN3700Telemetry summarises the current health of the index token.
type SYN3700Telemetry struct {
	Name            string    `json:"name"`
	Symbol          string    `json:"symbol"`
	ControllerCount int       `json:"controller_count"`
	ComponentCount  int       `json:"component_count"`
	TotalWeight     float64   `json:"total_weight"`
	LastRebalanced  time.Time `json:"last_rebalanced"`
	LastUpdated     time.Time `json:"last_updated"`
}

// SYN3700ControllerSnapshot is emitted in snapshots for persistence.
type SYN3700ControllerSnapshot struct {
	Address   Address   `json:"address"`
	GrantedBy Address   `json:"granted_by"`
	GrantedAt time.Time `json:"granted_at"`
}

// SYN3700ComponentSnapshot is exported for Stage 73 persistence.
type SYN3700ComponentSnapshot struct {
	Token          string    `json:"token"`
	Weight         float64   `json:"weight"`
	Drift          float64   `json:"drift"`
	Target         float64   `json:"target"`
	AddedBy        Address   `json:"added_by"`
	LastUpdated    time.Time `json:"last_updated"`
	LastRebalanced time.Time `json:"last_rebalanced"`
}

// SYN3700Snapshot encapsulates the full token state for Stage 73 persistence.
type SYN3700Snapshot struct {
	Name           string                      `json:"name"`
	Symbol         string                      `json:"symbol"`
	Components     []SYN3700ComponentSnapshot  `json:"components"`
	Controllers    []SYN3700ControllerSnapshot `json:"controllers"`
	Audit          []SYN3700AuditRecord        `json:"audit"`
	LastRebalanced time.Time                   `json:"last_rebalanced"`
	GeneratedAt    time.Time                   `json:"generated_at"`
}

// SYN3700Token models a weighted index token with controller governance.
type SYN3700Token struct {
	mu          sync.RWMutex
	name        string
	symbol      string
	components  map[string]syn3700Component
	controllers map[Address]syn3700Controller
	audit       []SYN3700AuditRecord
	lastUpdate  time.Time
	lastBalance time.Time
}

// NewSYN3700Token creates an empty index token.
func NewSYN3700Token(name, symbol string) *SYN3700Token {
	return &SYN3700Token{
		name:        name,
		symbol:      symbol,
		components:  make(map[string]syn3700Component),
		controllers: make(map[Address]syn3700Controller),
	}
}

func (t *SYN3700Token) appendAudit(actor Address, action, details string) {
	t.audit = append(t.audit, SYN3700AuditRecord{
		Timestamp: time.Now().UTC(),
		Actor:     actor,
		Action:    action,
		Details:   details,
	})
	if len(t.audit) > 256 {
		t.audit = append([]SYN3700AuditRecord(nil), t.audit[len(t.audit)-256:]...)
	}
}

func (t *SYN3700Token) controllerCountLocked() int { return len(t.controllers) }

func (t *SYN3700Token) isControllerLocked(addr Address) bool {
	_, ok := t.controllers[addr]
	return ok
}

// AddController grants governance rights to an address. The first controller can
// self-authorise to bootstrap the index. Subsequent additions require an existing
// controller to perform the operation.
func (t *SYN3700Token) AddController(addr, actor Address) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if addr == "" {
		return fmt.Errorf("controller address required")
	}
	if t.controllers == nil {
		t.controllers = make(map[Address]syn3700Controller)
	}
	if _, exists := t.controllers[addr]; exists {
		return fmt.Errorf("controller already exists")
	}
	if t.controllerCountLocked() > 0 && !t.isControllerLocked(actor) {
		return ErrControllerRequired
	}
	entry := syn3700Controller{Address: addr, GrantedBy: actor, GrantedAt: time.Now().UTC()}
	t.controllers[addr] = entry
	t.appendAudit(actor, "add_controller", string(addr))
	t.lastUpdate = time.Now().UTC()
	return nil
}

// RemoveController revokes governance rights, keeping at least one controller to
// avoid abandoning the token.
func (t *SYN3700Token) RemoveController(addr, actor Address) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.isControllerLocked(actor) {
		return ErrControllerRequired
	}
	if _, exists := t.controllers[addr]; !exists {
		return fmt.Errorf("controller not found")
	}
	if len(t.controllers) == 1 {
		return fmt.Errorf("cannot remove last controller")
	}
	delete(t.controllers, addr)
	t.appendAudit(actor, "remove_controller", string(addr))
	t.lastUpdate = time.Now().UTC()
	return nil
}

// Controllers returns the authorised controller list in deterministic order.
func (t *SYN3700Token) Controllers() []Address {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Address, 0, len(t.controllers))
	for addr := range t.controllers {
		out = append(out, addr)
	}
	sort.Slice(out, func(i, j int) bool { return string(out[i]) < string(out[j]) })
	return out
}

// IsController reports whether the provided address is authorised.
func (t *SYN3700Token) IsController(addr Address) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.isControllerLocked(addr)
}

// AddComponent registers a component token and its governance parameters.
func (t *SYN3700Token) AddComponent(token string, weight, drift float64, actor Address) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if token == "" {
		return fmt.Errorf("token required")
	}
	if !t.isControllerLocked(actor) {
		return ErrControllerRequired
	}
	if weight <= 0 {
		return fmt.Errorf("invalid weight")
	}
	if drift < 0 || drift > 1 {
		return fmt.Errorf("component drift must be between 0 and 1")
	}
	if _, exists := t.components[token]; exists {
		return ErrComponentExists
	}
	now := time.Now().UTC()
	t.components[token] = syn3700Component{
		Component: Component{
			Token:       token,
			Weight:      weight,
			Drift:       drift,
			Target:      weight,
			LastUpdated: now,
		},
		AddedBy: actor,
	}
	t.appendAudit(actor, "add_component", token)
	t.lastUpdate = now
	return nil
}

// UpdateComponent adjusts weight and drift tolerances without changing target.
func (t *SYN3700Token) UpdateComponent(token string, weight, drift float64, actor Address) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.isControllerLocked(actor) {
		return ErrControllerRequired
	}
	comp, ok := t.components[token]
	if !ok {
		return ErrComponentNotFound
	}
	if weight <= 0 {
		return fmt.Errorf("invalid weight")
	}
	if drift < 0 || drift > 1 {
		return fmt.Errorf("component drift must be between 0 and 1")
	}
	comp.Weight = weight
	comp.Drift = drift
	comp.Component.LastUpdated = time.Now().UTC()
	t.components[token] = comp
	t.appendAudit(actor, "update_component", token)
	t.lastUpdate = comp.Component.LastUpdated
	return nil
}

// RemoveComponent deletes a component from the index.
func (t *SYN3700Token) RemoveComponent(token string, actor Address) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.isControllerLocked(actor) {
		return ErrControllerRequired
	}
	if _, ok := t.components[token]; !ok {
		return ErrComponentNotFound
	}
	delete(t.components, token)
	t.appendAudit(actor, "remove_component", token)
	t.lastUpdate = time.Now().UTC()
	return nil
}

// ListComponents returns a snapshot of all components.
func (t *SYN3700Token) ListComponents() []Component {
	t.mu.RLock()
	defer t.mu.RUnlock()
	comps := make([]Component, 0, len(t.components))
	for _, comp := range t.components {
		comps = append(comps, comp.Component)
	}
	sort.Slice(comps, func(i, j int) bool { return comps[i].Token < comps[j].Token })
	return comps
}

// Value computes the index value given a map of token prices.
func (t *SYN3700Token) Value(prices map[string]float64) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var total float64
	for token, comp := range t.components {
		total += comp.Weight * prices[token]
	}
	return total
}

// RebalancePlan emits the current versus target weights for each component. The
// plan is always returned, even when current allocations already match targets,
// so downstream automation can verify drift budgets.
func (t *SYN3700Token) RebalancePlan() map[string][2]float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	plan := make(map[string][2]float64, len(t.components))
	for token, comp := range t.components {
		plan[token] = [2]float64{comp.Weight, comp.Target}
	}
	return plan
}

// RecordRebalance updates the target weight after a rebalance has been executed.
func (t *SYN3700Token) RecordRebalance(token string, newWeight float64, actor Address) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.isControllerLocked(actor) {
		return ErrControllerRequired
	}
	comp, ok := t.components[token]
	if !ok {
		return ErrComponentNotFound
	}
	if newWeight <= 0 {
		return fmt.Errorf("invalid weight")
	}
	comp.Target = newWeight
	comp.Weight = newWeight
	now := time.Now().UTC()
	comp.LastRebalanced = now
	comp.Component.LastUpdated = now
	t.components[token] = comp
	t.appendAudit(actor, "rebalance_component", fmt.Sprintf("%s:%f", token, newWeight))
	t.lastBalance = now
	t.lastUpdate = now
	return nil
}

// AuditLog returns a copy of recent audit events.
func (t *SYN3700Token) AuditLog() []SYN3700AuditRecord {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]SYN3700AuditRecord, len(t.audit))
	copy(out, t.audit)
	return out
}

// Telemetry summarises controllers, components and latest updates.
func (t *SYN3700Token) Telemetry() SYN3700Telemetry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var total float64
	var last time.Time
	for _, comp := range t.components {
		total += comp.Weight
		if comp.LastUpdated.After(last) {
			last = comp.LastUpdated
		}
	}
	if last.IsZero() {
		last = t.lastUpdate
	}
	return SYN3700Telemetry{
		Name:            t.name,
		Symbol:          t.symbol,
		ControllerCount: len(t.controllers),
		ComponentCount:  len(t.components),
		TotalWeight:     total,
		LastRebalanced:  t.lastBalance,
		LastUpdated:     last,
	}
}

// Snapshot captures the full token state for Stage 73 persistence.
func (t *SYN3700Token) Snapshot() SYN3700Snapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()
	components := make([]SYN3700ComponentSnapshot, 0, len(t.components))
	for _, comp := range t.components {
		components = append(components, SYN3700ComponentSnapshot{
			Token:          comp.Token,
			Weight:         comp.Weight,
			Drift:          comp.Drift,
			Target:         comp.Target,
			AddedBy:        comp.AddedBy,
			LastUpdated:    comp.LastUpdated,
			LastRebalanced: comp.LastRebalanced,
		})
	}
	sort.Slice(components, func(i, j int) bool { return components[i].Token < components[j].Token })

	controllers := make([]SYN3700ControllerSnapshot, 0, len(t.controllers))
	for _, ctrl := range t.controllers {
		controllers = append(controllers, SYN3700ControllerSnapshot{
			Address:   ctrl.Address,
			GrantedBy: ctrl.GrantedBy,
			GrantedAt: ctrl.GrantedAt,
		})
	}
	sort.Slice(controllers, func(i, j int) bool { return string(controllers[i].Address) < string(controllers[j].Address) })

	audit := make([]SYN3700AuditRecord, len(t.audit))
	copy(audit, t.audit)

	return SYN3700Snapshot{
		Name:           t.name,
		Symbol:         t.symbol,
		Components:     components,
		Controllers:    controllers,
		Audit:          audit,
		LastRebalanced: t.lastBalance,
		GeneratedAt:    time.Now().UTC(),
	}
}

// Restore rebuilds the token using a snapshot captured earlier.
func (t *SYN3700Token) Restore(snapshot SYN3700Snapshot) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.name = snapshot.Name
	t.symbol = snapshot.Symbol
	t.components = make(map[string]syn3700Component, len(snapshot.Components))
	for _, comp := range snapshot.Components {
		t.components[comp.Token] = syn3700Component{
			Component: Component{
				Token:       comp.Token,
				Weight:      comp.Weight,
				Drift:       comp.Drift,
				Target:      comp.Target,
				LastUpdated: comp.LastUpdated,
			},
			AddedBy:        comp.AddedBy,
			LastRebalanced: comp.LastRebalanced,
		}
	}
	t.controllers = make(map[Address]syn3700Controller, len(snapshot.Controllers))
	for _, ctrl := range snapshot.Controllers {
		t.controllers[ctrl.Address] = syn3700Controller{
			Address:   ctrl.Address,
			GrantedBy: ctrl.GrantedBy,
			GrantedAt: ctrl.GrantedAt,
		}
	}
	t.audit = append([]SYN3700AuditRecord(nil), snapshot.Audit...)
	t.lastBalance = snapshot.LastRebalanced
	t.lastUpdate = time.Now().UTC()
}

// DriftExceeded reports components exceeding their drift allowance.
func (t *SYN3700Token) DriftExceeded() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var exceeded []string
	for token, comp := range t.components {
		if math.Abs(comp.Weight-comp.Target) > comp.Drift {
			exceeded = append(exceeded, token)
		}
	}
	sort.Strings(exceeded)
	return exceeded
}
