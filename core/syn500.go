package core

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrInvalidGrant is returned when attempting to configure an invalid service tier.
var ErrInvalidGrant = errors.New("syn500: invalid grant")

// ErrUsageLimitReached indicates that a caller exhausted their allocation.
var ErrUsageLimitReached = errors.New("syn500: usage limit reached")

// ServiceTier defines access tiers for SYN500 utility tokens. Usage is tracked
// in units rather than single invocations so enterprise operators can enforce
// fine-grained quotas per address.
type ServiceTier struct {
	Tier int
	Max  uint64
	Used uint64
}

// UsageAudit records immutable metadata for a consumption event. The Digest
// field allows downstream systems to verify integrity when persisting the
// record to external stores.
type UsageAudit struct {
	Address   string
	Tier      int
	Amount    uint64
	Remaining uint64
	Note      string
	Timestamp time.Time
	Digest    string
}

// SYN500Token defines a tiered utility token with audit logging and
// concurrency-safe quota enforcement.
type SYN500Token struct {
	Name     string
	Symbol   string
	Owner    string
	Decimals uint8
	Supply   uint64

	mu     sync.RWMutex
	Grants map[string]*ServiceTier
	Audits []UsageAudit
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

// Grant assigns a service tier to an address. Existing grants are overwritten
// atomically so operators can raise or lower limits without dropping audit
// history.
func (t *SYN500Token) Grant(addr string, tier int, max uint64) error {
	if addr == "" || tier <= 0 || max == 0 {
		return ErrInvalidGrant
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if grant, ok := t.Grants[addr]; ok {
		grant.Tier = tier
		grant.Max = max
		if grant.Used > grant.Max {
			grant.Used = grant.Max
		}
		return nil
	}

	t.Grants[addr] = &ServiceTier{Tier: tier, Max: max}
	return nil
}

// Revoke removes the service tier for an address. It is safe to call multiple
// times.
func (t *SYN500Token) Revoke(addr string) {
	t.mu.Lock()
	delete(t.Grants, addr)
	t.mu.Unlock()
}

// Usage returns a snapshot of the current usage for an address.
func (t *SYN500Token) Usage(addr string) (ServiceTier, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	grant, ok := t.Grants[addr]
	if !ok {
		return ServiceTier{}, false
	}
	return *grant, true
}

// Use records a token usage for an address. The amount must be positive and
// within the caller's remaining allocation.
func (t *SYN500Token) Use(addr string, amount uint64, note string) (UsageAudit, error) {
	if amount == 0 {
		return UsageAudit{}, fmt.Errorf("syn500: amount must be positive")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	grant, ok := t.Grants[addr]
	if !ok {
		return UsageAudit{}, errors.New("syn500: no tier granted")
	}
	if grant.Used+amount > grant.Max {
		return UsageAudit{}, ErrUsageLimitReached
	}

	grant.Used += amount
	record := UsageAudit{
		Address:   addr,
		Tier:      grant.Tier,
		Amount:    amount,
		Remaining: grant.Max - grant.Used,
		Note:      note,
		Timestamp: time.Now().UTC(),
	}
	digest := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%d|%s|%d", record.Address, record.Tier, record.Amount, record.Note, record.Timestamp.UnixNano())))
	record.Digest = hex.EncodeToString(digest[:])
	t.Audits = append(t.Audits, record)
	return record, nil
}

// Snapshot returns a copy of the current grants and audit length for external
// monitoring.
func (t *SYN500Token) Snapshot() map[string]any {
	t.mu.RLock()
	defer t.mu.RUnlock()

	grants := make(map[string]ServiceTier, len(t.Grants))
	for addr, grant := range t.Grants {
		grants[addr] = *grant
	}
	return map[string]any{
		"name":     t.Name,
		"symbol":   t.Symbol,
		"owner":    t.Owner,
		"grants":   grants,
		"audits":   len(t.Audits),
		"supply":   t.Supply,
		"decimals": t.Decimals,
	}
}

// AuditTrail returns the most recent usage records up to limit entries. If
// limit is zero the full audit trail is returned.
func (t *SYN500Token) AuditTrail(limit int) []UsageAudit {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if limit <= 0 || limit >= len(t.Audits) {
		out := make([]UsageAudit, len(t.Audits))
		copy(out, t.Audits)
		return out
	}
	start := len(t.Audits) - limit
	out := make([]UsageAudit, len(t.Audits[start:]))
	copy(out, t.Audits[start:])
	return out
}
