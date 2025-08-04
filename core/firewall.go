package core

import "sync"

// Firewall manages simple allow and block lists for peer IP addresses.  It is
// designed for lightweight in-memory use by node implementations.
type Firewall struct {
	mu      sync.RWMutex
	allowed map[string]struct{}
	blocked map[string]struct{}
}

// NewFirewall returns an initialised Firewall instance.
func NewFirewall() *Firewall {
	return &Firewall{
		allowed: make(map[string]struct{}),
		blocked: make(map[string]struct{}),
	}
}

// Allow permits connections from the given IP address and removes any explicit
// block rule that may exist.
func (f *Firewall) Allow(ip string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.blocked, ip)
	f.allowed[ip] = struct{}{}
}

// Block denies connections from the given IP address and removes any existing
// allow rule.
func (f *Firewall) Block(ip string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.allowed, ip)
	f.blocked[ip] = struct{}{}
}

// IsAllowed returns true if the given IP address is permitted by the firewall.
// If no allow rules are defined, the firewall defaults to allowing all addresses
// not explicitly blocked.
func (f *Firewall) IsAllowed(ip string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if _, banned := f.blocked[ip]; banned {
		return false
	}
	if len(f.allowed) == 0 {
		return true
	}
	_, ok := f.allowed[ip]
	return ok
}

// Rules returns the current allow and block lists for inspection.
func (f *Firewall) Rules() (allowed []string, blocked []string) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for ip := range f.allowed {
		allowed = append(allowed, ip)
	}
	for ip := range f.blocked {
		blocked = append(blocked, ip)
	}
	return
}
