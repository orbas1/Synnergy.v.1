package core

import (
	"errors"
	"sync"
	"time"
)

// ServiceTier defines access tiers for SYN500 utility tokens.
type ServiceTier struct {
	Tier          int           `json:"tier"`
	Max           uint64        `json:"max"`
	Used          uint64        `json:"used"`
	Window        time.Duration `json:"-"`
	WindowSeconds int64         `json:"window_seconds"`
	LastReset     time.Time     `json:"-"`
}

// SYN500Token defines a utility token with usage tracking and grant telemetry.
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
	return &SYN500Token{
		Name:     name,
		Symbol:   symbol,
		Owner:    owner,
		Decimals: decimals,
		Supply:   supply,
		Grants:   make(map[string]*ServiceTier),
		Created:  time.Now().UTC(),
	}
}

// Grant assigns a service tier to an address.
func (t *SYN500Token) Grant(addr string, tier int, max uint64, window time.Duration) error {
	if addr == "" {
		return errors.New("address required")
	}
	if tier <= 0 {
		return errors.New("tier must be positive")
	}
	if max == 0 {
		return errors.New("max usage must be positive")
	}
	if window <= 0 {
		window = time.Hour
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Grants[addr] = &ServiceTier{
		Tier:          tier,
		Max:           max,
		Used:          0,
		Window:        window,
		WindowSeconds: int64(window.Seconds()),
		LastReset:     time.Now().UTC(),
	}
	return nil
}

// Use records a token usage for an address.
func (t *SYN500Token) Use(addr string, now time.Time) error {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	g, ok := t.Grants[addr]
	if !ok {
		return errors.New("no tier granted")
	}
	if g.Window > 0 && now.Sub(g.LastReset) >= g.Window {
		g.Used = 0
		g.LastReset = now
	}
	if g.Used >= g.Max {
		return errors.New("usage limit reached")
	}
	g.Used++
	return nil
}

// Status returns a snapshot of the grant for addr.
func (t *SYN500Token) Status(addr string) (ServiceTier, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	g, ok := t.Grants[addr]
	if !ok {
		return ServiceTier{}, false
	}
	cp := *g
	return ServiceTier{
		Tier:          cp.Tier,
		Max:           cp.Max,
		Used:          cp.Used,
		WindowSeconds: cp.WindowSeconds,
	}, true
}

// Telemetry summarises usage grants for monitoring.
type SYN500Telemetry struct {
	Grants int `json:"grants"`
}

// Telemetry reports basic usage metadata.
func (t *SYN500Token) Telemetry() SYN500Telemetry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return SYN500Telemetry{Grants: len(t.Grants)}
}
