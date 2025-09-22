package tokens

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// SYN20Token extends BaseToken with pause and freeze capabilities.
type SYN20Token struct {
	*BaseToken
	mu               sync.RWMutex
	paused           bool
	frozen           map[string]bool
	whitelist        map[string]bool
	enforceWhitelist bool
	circuitBreaker   bool
	limits           map[string]uint64
	spent            map[string]uint64
	windowStart      map[string]time.Time
	audit            []SYN20AuditRecord
}

// NewSYN20Token creates a new SYN20 token.
func NewSYN20Token(id TokenID, name, symbol string, decimals uint8) *SYN20Token {
	return &SYN20Token{
		BaseToken:   NewBaseToken(id, name, symbol, decimals),
		frozen:      make(map[string]bool),
		whitelist:   make(map[string]bool),
		limits:      make(map[string]uint64),
		spent:       make(map[string]uint64),
		windowStart: make(map[string]time.Time),
	}
}

// Pause halts all transfer, mint and burn operations.
func (t *SYN20Token) Pause() {
	t.mu.Lock()
	t.paused = true
	t.mu.Unlock()
}

// Unpause resumes operations.
func (t *SYN20Token) Unpause() {
	t.mu.Lock()
	t.paused = false
	t.mu.Unlock()
}

// Freeze prevents an address from participating in transfers.
func (t *SYN20Token) Freeze(addr string) {
	t.mu.Lock()
	t.frozen[addr] = true
	t.mu.Unlock()
}

// Unfreeze lifts restrictions on an address.
func (t *SYN20Token) Unfreeze(addr string) {
	t.mu.Lock()
	delete(t.frozen, addr)
	t.mu.Unlock()
}

// SetWhitelistRequirement toggles whitelist enforcement for transfers.
func (t *SYN20Token) SetWhitelistRequirement(enabled bool) {
	t.mu.Lock()
	t.enforceWhitelist = enabled
	t.mu.Unlock()
}

// AddToWhitelist grants access to an address.
func (t *SYN20Token) AddToWhitelist(addr string) {
	t.mu.Lock()
	t.whitelist[addr] = true
	t.mu.Unlock()
}

// RemoveFromWhitelist removes a whitelisted address.
func (t *SYN20Token) RemoveFromWhitelist(addr string) {
	t.mu.Lock()
	delete(t.whitelist, addr)
	t.mu.Unlock()
}

// SetCircuitBreaker toggles the enterprise circuit breaker.
func (t *SYN20Token) SetCircuitBreaker(active bool) {
	t.mu.Lock()
	t.circuitBreaker = active
	t.mu.Unlock()
}

// SetDailyLimit configures a rolling 24h limit for an address.
func (t *SYN20Token) SetDailyLimit(addr string, amount uint64) {
	t.mu.Lock()
	t.limits[addr] = amount
	if _, ok := t.windowStart[addr]; !ok {
		t.windowStart[addr] = time.Now()
	}
	t.mu.Unlock()
}

// SYN20AuditRecord records privileged state transitions for monitoring.
type SYN20AuditRecord struct {
	From      string
	To        string
	Amount    uint64
	Reason    string
	Timestamp time.Time
}

var (
	// ErrSYN20CircuitBreakerActive indicates transfers halted network wide.
	ErrSYN20CircuitBreakerActive = errors.New("syn20 circuit breaker active")
	// ErrSYN20WhitelistViolation indicates missing whitelist permission.
	ErrSYN20WhitelistViolation = errors.New("recipient not whitelisted")
	// ErrSYN20LimitExceeded indicates daily rolling limit exceeded.
	ErrSYN20LimitExceeded = errors.New("daily limit exceeded")
)

// Transfer overrides BaseToken.Transfer adding pause/freeze checks.
func (t *SYN20Token) Transfer(from, to string, amount uint64) error {
	t.mu.RLock()
	paused := t.paused
	breaker := t.circuitBreaker
	fromFrozen := t.frozen[from]
	toFrozen := t.frozen[to]
	enforce := t.enforceWhitelist
	toAllowed := t.whitelist[to]
	limit := t.limits[from]
	spent := t.spent[from]
	start := t.windowStart[from]
	t.mu.RUnlock()
	if paused {
		return fmt.Errorf("token transfers are paused")
	}
	if breaker {
		return ErrSYN20CircuitBreakerActive
	}
	if fromFrozen || toFrozen {
		return fmt.Errorf("address frozen")
	}
	if enforce && !toAllowed {
		return ErrSYN20WhitelistViolation
	}
	if start.IsZero() {
		start = time.Now()
	}
	if time.Since(start) >= 24*time.Hour {
		t.mu.Lock()
		t.spent[from] = 0
		t.windowStart[from] = time.Now()
		spent = 0
		t.mu.Unlock()
	}
	if limit > 0 && spent+amount > limit {
		return ErrSYN20LimitExceeded
	}
	return t.BaseToken.Transfer(from, to, amount)
}

// Mint creates tokens if operations are not paused or frozen.
func (t *SYN20Token) Mint(to string, amount uint64) error {
	t.mu.RLock()
	paused := t.paused
	frozen := t.frozen[to]
	t.mu.RUnlock()
	if paused || frozen {
		return fmt.Errorf("minting restricted")
	}
	return t.BaseToken.Mint(to, amount)
}

// Burn destroys tokens if operations are allowed.
func (t *SYN20Token) Burn(from string, amount uint64) error {
	t.mu.RLock()
	paused := t.paused
	frozen := t.frozen[from]
	t.mu.RUnlock()
	if paused || frozen {
		return fmt.Errorf("burning restricted")
	}
	return t.BaseToken.Burn(from, amount)
}

// RecordAudit appends a traceable audit record for compliance reporting.
func (t *SYN20Token) RecordAudit(from, to string, amount uint64, reason string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.audit = append(t.audit, SYN20AuditRecord{From: from, To: to, Amount: amount, Reason: reason, Timestamp: time.Now()})
	t.spent[from] += amount
	if _, ok := t.windowStart[from]; !ok {
		t.windowStart[from] = time.Now()
	}
}

// AuditTrail returns the most recent audit entries.
func (t *SYN20Token) AuditTrail(limit int) []SYN20AuditRecord {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if limit <= 0 || limit >= len(t.audit) {
		out := make([]SYN20AuditRecord, len(t.audit))
		copy(out, t.audit)
		return out
	}
	out := make([]SYN20AuditRecord, limit)
	copy(out, t.audit[len(t.audit)-limit:])
	return out
}
