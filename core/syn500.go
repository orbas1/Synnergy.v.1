package core

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// UsageEvent captures state transitions for auditability and telemetry.
type UsageEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Details   string    `json:"details,omitempty"`
}

// ServiceTier defines access tiers for SYN500 utility tokens and records usage.
type ServiceTier struct {
	Tier        int           `json:"tier"`
	Max         uint64        `json:"max"`
	Used        uint64        `json:"used"`
	Window      time.Duration `json:"-"`
	LastReset   time.Time     `json:"last_reset,omitempty"`
	AuditTrail  []UsageEvent  `json:"audit,omitempty"`
	windowCache string        `json:"window,omitempty"`
}

// clone returns a deep copy suitable for external consumers.
func (s *ServiceTier) clone() ServiceTier {
	if s == nil {
		return ServiceTier{}
	}
	out := *s
	if len(s.AuditTrail) > 0 {
		out.AuditTrail = append([]UsageEvent(nil), s.AuditTrail...)
	}
	if s.Window > 0 {
		out.windowCache = s.Window.String()
	}
	return out
}

// SYN500Telemetry summarises grant utilisation for monitoring endpoints.
type SYN500Telemetry struct {
	Grants      int     `json:"grants"`
	Windowed    int     `json:"windowed"`
	TotalUsage  uint64  `json:"total_usage"`
	Utilisation float64 `json:"utilisation"`
}

// SYN500Token defines a concurrency-safe utility token with usage tracking.
type SYN500Token struct {
	mu        sync.RWMutex
	Name      string
	Symbol    string
	Owner     string
	Decimals  uint8
	Supply    uint64
	CreatedAt time.Time
	Grants    map[string]*ServiceTier
	Audit     []UsageEvent
}

// NewSYN500Token creates a new utility token.
func NewSYN500Token(name, symbol, owner string, decimals uint8, supply uint64) *SYN500Token {
	return &SYN500Token{
		Name:      name,
		Symbol:    symbol,
		Owner:     owner,
		Decimals:  decimals,
		Supply:    supply,
		CreatedAt: time.Now().UTC(),
		Grants:    make(map[string]*ServiceTier),
	}
}

// Grant assigns or updates a service tier to an address.
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
	if window < 0 {
		return errors.New("window cannot be negative")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.Grants == nil {
		t.Grants = make(map[string]*ServiceTier)
	}
	now := time.Now().UTC()
	tierState := &ServiceTier{Tier: tier, Max: max, Window: window, LastReset: now}
	if existing, ok := t.Grants[addr]; ok {
		tierState.Used = existing.Used
		if existing.Window == window {
			tierState.LastReset = existing.LastReset
		}
		tierState.AuditTrail = append(existing.AuditTrail, UsageEvent{Timestamp: now, Actor: addr, Action: "updated", Details: fmt.Sprintf("tier=%d max=%d", tier, max)})
	} else {
		tierState.AuditTrail = []UsageEvent{{Timestamp: now, Actor: addr, Action: "granted", Details: fmt.Sprintf("tier=%d max=%d", tier, max)}}
	}
	if window == 0 {
		tierState.LastReset = time.Time{}
	}
	t.Grants[addr] = tierState
	t.Audit = append(t.Audit, UsageEvent{Timestamp: now, Actor: addr, Action: "grant", Details: fmt.Sprintf("tier=%d max=%d", tier, max)})
	return nil
}

// Use records a token usage for an address.
func (t *SYN500Token) Use(addr string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	g, ok := t.Grants[addr]
	if !ok {
		return errors.New("no tier granted")
	}
	now := time.Now().UTC()
	if g.Window > 0 && !g.LastReset.IsZero() && now.Sub(g.LastReset) >= g.Window {
		g.Used = 0
		g.LastReset = now
		g.AuditTrail = append(g.AuditTrail, UsageEvent{Timestamp: now, Actor: addr, Action: "window_reset", Details: g.Window.String()})
	}
	if g.Used >= g.Max {
		return errors.New("usage limit reached")
	}
	g.Used++
	g.AuditTrail = append(g.AuditTrail, UsageEvent{Timestamp: now, Actor: addr, Action: "use", Details: fmt.Sprintf("usage=%d", g.Used)})
	t.Audit = append(t.Audit, UsageEvent{Timestamp: now, Actor: addr, Action: "use"})
	return nil
}

// Status returns the current service tier details for an address.
func (t *SYN500Token) Status(addr string) (ServiceTier, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	tier, ok := t.Grants[addr]
	if !ok {
		return ServiceTier{}, false
	}
	clone := tier.clone()
	return clone, true
}

// Telemetry summarises utilisation across all grants.
func (t *SYN500Token) Telemetry() SYN500Telemetry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	telemetry := SYN500Telemetry{}
	telemetry.Grants = len(t.Grants)
	if telemetry.Grants == 0 {
		return telemetry
	}
	for _, tier := range t.Grants {
		telemetry.TotalUsage += tier.Used
		if tier.Window > 0 {
			telemetry.Windowed++
		}
		if tier.Max > 0 {
			telemetry.Utilisation += float64(tier.Used) / float64(tier.Max)
		}
	}
	telemetry.Utilisation = telemetry.Utilisation / float64(telemetry.Grants)
	return telemetry
}

// AuditLog returns a copy of historical events for compliance surfaces.
func (t *SYN500Token) AuditLog() []UsageEvent {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]UsageEvent, len(t.Audit))
	copy(out, t.Audit)
	return out
}
