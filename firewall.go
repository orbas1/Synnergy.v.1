package synnergy

import "sync"

// Firewall manages block lists for wallet addresses, token identifiers and
// peer IP addresses. It is intentionally lightweight and safe for concurrent
// use by multiple goroutines.
type Firewall struct {
	mu            sync.RWMutex
	blockedAddrs  map[string]struct{}
	blockedTokens map[string]struct{}
	blockedIPs    map[string]struct{}
}

// NewFirewall creates an empty Firewall instance.
func NewFirewall() *Firewall {
	return &Firewall{
		blockedAddrs:  make(map[string]struct{}),
		blockedTokens: make(map[string]struct{}),
		blockedIPs:    make(map[string]struct{}),
	}
}

// BlockAddress adds an address to the block list.
func (f *Firewall) BlockAddress(addr string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.blockedAddrs[addr] = struct{}{}
}

// UnblockAddress removes an address from the block list.
func (f *Firewall) UnblockAddress(addr string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.blockedAddrs, addr)
}

// IsAddressBlocked returns true if the address is currently blocked.
func (f *Firewall) IsAddressBlocked(addr string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	_, ok := f.blockedAddrs[addr]
	return ok
}

// BlockToken adds a token identifier to the block list.
func (f *Firewall) BlockToken(id string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.blockedTokens[id] = struct{}{}
}

// UnblockToken removes a token identifier from the block list.
func (f *Firewall) UnblockToken(id string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.blockedTokens, id)
}

// IsTokenBlocked returns true if the token identifier is blocked.
func (f *Firewall) IsTokenBlocked(id string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	_, ok := f.blockedTokens[id]
	return ok
}

// BlockIP adds a peer IP address to the block list.
func (f *Firewall) BlockIP(ip string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.blockedIPs[ip] = struct{}{}
}

// UnblockIP removes a peer IP address from the block list.
func (f *Firewall) UnblockIP(ip string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.blockedIPs, ip)
}

// IsIPBlocked returns true if the IP address is blocked.
func (f *Firewall) IsIPBlocked(ip string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	_, ok := f.blockedIPs[ip]
	return ok
}

// Rules returns slices of currently blocked addresses, tokens and IPs.
func (f *Firewall) Rules() (addrs, tokens, ips []string) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for a := range f.blockedAddrs {
		addrs = append(addrs, a)
	}
	for t := range f.blockedTokens {
		tokens = append(tokens, t)
	}
	for ip := range f.blockedIPs {
		ips = append(ips, ip)
	}
	return
}
