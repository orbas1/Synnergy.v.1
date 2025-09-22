package core

import (
	"sync"
	"time"
)

// Firewall manages allow and block lists for peer IP addresses with optional
// expiry and audit metadata. It is designed for lightweight in-memory use by
// node implementations.
type Firewall struct {
	mu           sync.RWMutex
	allowed      map[string]ruleRecord
	blocked      map[string]ruleRecord
	defaultAllow bool
}

type ruleRecord struct {
	reason     string
	created    time.Time
	expires    time.Time
	hits       uint64
	persistent bool
}

// RuleOptions define metadata stored alongside firewall entries.
type RuleOptions struct {
	Reason     string
	TTL        time.Duration
	Persistent bool
}

// NewFirewall returns an initialised Firewall instance.
func NewFirewall() *Firewall {
	return &Firewall{
		allowed:      make(map[string]ruleRecord),
		blocked:      make(map[string]ruleRecord),
		defaultAllow: true,
	}
}

// Allow permits connections from the given IP address and removes any explicit
// block rule that may exist.
func (f *Firewall) Allow(ip string) {
	f.AllowWithOptions(ip, RuleOptions{})
}

// AllowWithOptions registers an allow rule with metadata.
func (f *Firewall) AllowWithOptions(ip string, opts RuleOptions) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.blocked, ip)
	f.allowed[ip] = ruleRecord{
		reason:     opts.Reason,
		created:    time.Now(),
		expires:    expiryFromOptions(opts),
		persistent: opts.Persistent,
	}
}

// Block denies connections from the given IP address and removes any existing
// allow rule.
func (f *Firewall) Block(ip string) {
	f.BlockWithOptions(ip, RuleOptions{})
}

// BlockWithOptions registers a block rule with metadata.
func (f *Firewall) BlockWithOptions(ip string, opts RuleOptions) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.allowed, ip)
	f.blocked[ip] = ruleRecord{
		reason:     opts.Reason,
		created:    time.Now(),
		expires:    expiryFromOptions(opts),
		persistent: opts.Persistent,
	}
}

// IsAllowed returns true if the given IP address is permitted by the firewall.
// If no allow rules are defined, the firewall defaults to allowing all addresses
// not explicitly blocked.
func (f *Firewall) IsAllowed(ip string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.pruneLocked(time.Now())
	if rec, banned := f.blocked[ip]; banned {
		rec.hits++
		f.blocked[ip] = rec
		return false
	}
	if len(f.allowed) == 0 {
		return f.defaultAllow
	}
	rec, ok := f.allowed[ip]
	if ok {
		rec.hits++
		f.allowed[ip] = rec
	}
	return ok
}

// Rules returns the current allow and block lists for inspection.
func (f *Firewall) Rules() (allowed []string, blocked []string) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for ip, rec := range f.allowed {
		if !rec.persistent && expired(rec.expires, time.Now()) {
			continue
		}
		allowed = append(allowed, ip)
	}
	for ip, rec := range f.blocked {
		if !rec.persistent && expired(rec.expires, time.Now()) {
			continue
		}
		blocked = append(blocked, ip)
	}
	return
}

// Reset clears all allow and block rules.
func (f *Firewall) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.allowed = make(map[string]ruleRecord)
	f.blocked = make(map[string]ruleRecord)
	f.defaultAllow = true
}

// SetDefaultAllow toggles whether addresses without explicit allow rules are
// permitted.
func (f *Firewall) SetDefaultAllow(allow bool) {
	f.mu.Lock()
	f.defaultAllow = allow
	f.mu.Unlock()
}

// PruneExpired removes expired rules and returns the count of deleted entries.
func (f *Firewall) PruneExpired() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.pruneLocked(time.Now())
}

// RuleDetail exposes metadata about a firewall rule for observability.
type RuleDetail struct {
	IP         string
	Reason     string
	Created    time.Time
	Expires    time.Time
	Hits       uint64
	Allowed    bool
	Persistent bool
}

// RuleDetails returns structured information about allow and block entries.
func (f *Firewall) RuleDetails() []RuleDetail {
	now := time.Now()
	f.mu.RLock()
	defer f.mu.RUnlock()
	var out []RuleDetail
	for ip, rec := range f.allowed {
		if !rec.persistent && expired(rec.expires, now) {
			continue
		}
		out = append(out, RuleDetail{
			IP:         ip,
			Reason:     rec.reason,
			Created:    rec.created,
			Expires:    rec.expires,
			Hits:       rec.hits,
			Allowed:    true,
			Persistent: rec.persistent,
		})
	}
	for ip, rec := range f.blocked {
		if !rec.persistent && expired(rec.expires, now) {
			continue
		}
		out = append(out, RuleDetail{
			IP:         ip,
			Reason:     rec.reason,
			Created:    rec.created,
			Expires:    rec.expires,
			Hits:       rec.hits,
			Allowed:    false,
			Persistent: rec.persistent,
		})
	}
	return out
}

func expiryFromOptions(opts RuleOptions) time.Time {
	if opts.Persistent || opts.TTL <= 0 {
		return time.Time{}
	}
	return time.Now().Add(opts.TTL)
}

func expired(expires time.Time, now time.Time) bool {
	if expires.IsZero() {
		return false
	}
	return now.After(expires)
}

func (f *Firewall) pruneLocked(now time.Time) int {
	removed := 0
	for ip, rec := range f.allowed {
		if rec.persistent || !expired(rec.expires, now) {
			continue
		}
		delete(f.allowed, ip)
		removed++
	}
	for ip, rec := range f.blocked {
		if rec.persistent || !expired(rec.expires, now) {
			continue
		}
		delete(f.blocked, ip)
		removed++
	}
	return removed
}
