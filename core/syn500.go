package core

import (
	"errors"
	"sync"
	"time"
)

// ServiceTier defines access tiers for SYN500 utility tokens.
type ServiceTier struct {
	Tier      int           `json:"tier"`
	Max       uint64        `json:"max"`
	Used      uint64        `json:"used"`
	Window    time.Duration `json:"window"`
	LastReset time.Time     `json:"last_reset"`
}

// SYN500Token defines a utility token with usage tracking and telemetry.
type SYN500Token struct {
	mu       sync.RWMutex
	Name     string
	Symbol   string
	Owner    string
	Decimals uint8
	Supply   uint64
	Grants   map[string]*ServiceTier
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
	}
}

// Grant assigns a service tier to an address with a usage window.
func (t *SYN500Token) Grant(addr string, tier int, max uint64, window time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if window <= 0 {
		window = time.Hour
	}
	t.Grants[addr] = &ServiceTier{Tier: tier, Max: max, Window: window, LastReset: time.Now().UTC()}
}

func (t *SYN500Token) resetIfNeeded(st *ServiceTier) {
	if st.Window <= 0 {
		st.Window = time.Hour
	}
	if time.Since(st.LastReset) >= st.Window {
		st.Used = 0
		st.LastReset = time.Now().UTC()
	}
}

// Use records a token usage for an address.
func (t *SYN500Token) Use(addr string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	st, ok := t.Grants[addr]
	if !ok {
		return errors.New("no tier granted")
	}
	t.resetIfNeeded(st)
	if st.Used >= st.Max {
		return errors.New("usage limit reached")
	}
	st.Used++
	return nil
}

// Status returns a snapshot of the service tier for an address.
func (t *SYN500Token) Status(addr string) (ServiceTier, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	st, ok := t.Grants[addr]
	if !ok {
		return ServiceTier{}, false
	}
	cp := *st
	return cp, true
}

// UtilityTelemetry summarises grant distribution for monitoring.
type UtilityTelemetry struct {
	Grants int `json:"grants"`
	Active int `json:"active"`
}

// Telemetry returns aggregated usage information.
func (t *SYN500Token) Telemetry() UtilityTelemetry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var tele UtilityTelemetry
	tele.Grants = len(t.Grants)
	for _, st := range t.Grants {
		if st.Used < st.Max {
			tele.Active++
		}
	}
	return tele
}
