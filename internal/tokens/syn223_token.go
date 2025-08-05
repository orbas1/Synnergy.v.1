package tokens

import (
	"errors"
	"sync"
)

// SYN223Token implements a secure transfer token with whitelist and blacklist controls.
type SYN223Token struct {
	mu        sync.RWMutex
	Name      string
	Symbol    string
	Owner     string
	balances  map[string]uint64
	whitelist map[string]bool
	blacklist map[string]bool
}

// NewSYN223Token creates a new SYN223 token and assigns the initial supply to the owner.
func NewSYN223Token(name, symbol, owner string, supply uint64) *SYN223Token {
	t := &SYN223Token{
		Name:      name,
		Symbol:    symbol,
		Owner:     owner,
		balances:  map[string]uint64{owner: supply},
		whitelist: make(map[string]bool),
		blacklist: make(map[string]bool),
	}
	t.whitelist[owner] = true
	return t
}

// AddToWhitelist authorises an address to receive tokens.
func (t *SYN223Token) AddToWhitelist(addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.whitelist[addr] = true
}

// RemoveFromWhitelist removes an address from the whitelist.
func (t *SYN223Token) RemoveFromWhitelist(addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.whitelist, addr)
}

// AddToBlacklist blocks an address from participating in transfers.
func (t *SYN223Token) AddToBlacklist(addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.blacklist[addr] = true
}

// RemoveFromBlacklist lifts a previously applied blacklist restriction.
func (t *SYN223Token) RemoveFromBlacklist(addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.blacklist, addr)
}

// Transfer performs a safe token transfer verifying whitelist and blacklist rules.
func (t *SYN223Token) Transfer(from, to string, amount uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.blacklist[from] || t.blacklist[to] {
		return errors.New("address blacklisted")
	}
	if !t.whitelist[to] {
		return errors.New("recipient not whitelisted")
	}
	bal := t.balances[from]
	if bal < amount {
		return errors.New("insufficient balance")
	}
	t.balances[from] = bal - amount
	t.balances[to] += amount
	return nil
}

// BalanceOf returns the current balance of an address.
func (t *SYN223Token) BalanceOf(addr string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.balances[addr]
}
