package core

import (
	"errors"
	"sync"
	"time"
)

// ServiceTier defines access tiers for SYN500 utility tokens.
type ServiceTier struct {
	Tier         int           `json:"tier"`
	Max          uint64        `json:"max"`
	Used         uint64        `json:"used"`
	Window       time.Duration `json:"window"`
	LastReset    time.Time     `json:"last_reset"`
	TotalGranted uint64        `json:"total_granted"`
}

// copy returns a defensive clone of the service tier.
func (s *ServiceTier) copy() *ServiceTier {
	if s == nil {
		return nil
	}
	cp := *s
	return &cp
}

// SYN500Token defines a utility token with usage tracking and rate limiting windows.
type SYN500Token struct {
	mu       sync.RWMutex
	Name     string                  `json:"name"`
	Symbol   string                  `json:"symbol"`
	Owner    string                  `json:"owner"`
	Decimals uint8                   `json:"decimals"`
	Supply   uint64                  `json:"supply"`
	Grants   map[string]*ServiceTier `json:"grants"`
	Created  time.Time               `json:"created_at"`
}

// NewSYN500Token creates a new utility token.
func NewSYN500Token(name, symbol, owner string, decimals uint8, supply uint64) *SYN500Token {
	return &SYN500Token{Name: name, Symbol: symbol, Owner: owner, Decimals: decimals, Supply: supply, Grants: make(map[string]*ServiceTier), Created: time.Now().UTC()}
}

// RestoreSYN500Token rebuilds a token from a persisted snapshot.
func RestoreSYN500Token(snapshot *SYN500Snapshot) *SYN500Token {
	if snapshot == nil {
		return nil
	}
	token := &SYN500Token{
		Name:     snapshot.Name,
		Symbol:   snapshot.Symbol,
		Owner:    snapshot.Owner,
		Decimals: snapshot.Decimals,
		Supply:   snapshot.Supply,
		Grants:   make(map[string]*ServiceTier),
		Created:  snapshot.Created,
	}
	for addr, tier := range snapshot.Grants {
		token.Grants[addr] = tier.copy()
	}
	return token
}

// Grant assigns a service tier to an address.
func (t *SYN500Token) Grant(addr string, tier int, max uint64, window time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.Grants == nil {
		t.Grants = make(map[string]*ServiceTier)
	}
	t.Grants[addr] = &ServiceTier{Tier: tier, Max: max, Window: window, LastReset: time.Now().UTC()}
}

// UseAt records a token usage for an address at the given time.
func (t *SYN500Token) UseAt(addr string, now time.Time) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	g, ok := t.Grants[addr]
	if !ok {
		return errors.New("no tier granted")
	}
	if g.Window > 0 && !g.LastReset.IsZero() && now.Sub(g.LastReset) >= g.Window {
		g.Used = 0
		g.LastReset = now
	}
	if g.Used >= g.Max {
		return errors.New("usage limit reached")
	}
	g.Used++
	g.TotalGranted++
	if g.LastReset.IsZero() {
		g.LastReset = now
	}
	return nil
}

// Use records a token usage for an address with the current time.
func (t *SYN500Token) Use(addr string) error {
	return t.UseAt(addr, time.Now().UTC())
}

// Status returns usage information for an address.
func (t *SYN500Token) Status(addr string) (*ServiceTier, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	tier, ok := t.Grants[addr]
	if !ok {
		return nil, false
	}
	return tier.copy(), true
}

// Telemetry summarises active grants for dashboards.
func (t *SYN500Token) Telemetry() SYN500Telemetry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	tele := SYN500Telemetry{}
	tele.TotalGrants = len(t.Grants)
	for addr, tier := range t.Grants {
		if tier == nil {
			continue
		}
		tele.TotalUsage += tier.Used
		if tier.Used < tier.Max {
			tele.ActiveAddresses = append(tele.ActiveAddresses, addr)
		}
	}
	return tele
}

// Snapshot returns a serialisable view of the token.
func (t *SYN500Token) Snapshot() *SYN500Snapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()
	snap := &SYN500Snapshot{
		Name:     t.Name,
		Symbol:   t.Symbol,
		Owner:    t.Owner,
		Decimals: t.Decimals,
		Supply:   t.Supply,
		Created:  t.Created,
		Grants:   make(map[string]*ServiceTier),
	}
	for addr, tier := range t.Grants {
		snap.Grants[addr] = tier.copy()
	}
	return snap
}

// SYN500Telemetry aggregates token activity for reporting.
type SYN500Telemetry struct {
	TotalGrants     int      `json:"grants"`
	TotalUsage      uint64   `json:"usage"`
	ActiveAddresses []string `json:"active_addresses"`
}

// SYN500Snapshot captures the persisted state of the token.
type SYN500Snapshot struct {
	Name     string                  `json:"name"`
	Symbol   string                  `json:"symbol"`
	Owner    string                  `json:"owner"`
	Decimals uint8                   `json:"decimals"`
	Supply   uint64                  `json:"supply"`
	Grants   map[string]*ServiceTier `json:"grants"`
	Created  time.Time               `json:"created_at"`
}
