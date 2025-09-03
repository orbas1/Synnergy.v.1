package core

import (
	"errors"
	"sync"
)

// DAOTokenLedger tracks DAO membership token balances.
type DAOTokenLedger struct {
	mu       sync.RWMutex
	balances map[string]uint64
}

// NewDAOTokenLedger returns an initialised ledger.
func NewDAOTokenLedger() *DAOTokenLedger {
	return &DAOTokenLedger{balances: make(map[string]uint64)}
}

// Mint creates tokens for an address.
func (l *DAOTokenLedger) Mint(addr string, amount uint64) {
	l.mu.Lock()
	l.balances[addr] += amount
	l.mu.Unlock()
}

// Transfer moves tokens between addresses.
func (l *DAOTokenLedger) Transfer(from, to string, amount uint64) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.balances[from] < amount {
		return errors.New("insufficient balance")
	}
	l.balances[from] -= amount
	l.balances[to] += amount
	return nil
}

// Balance returns the token balance for an address.
func (l *DAOTokenLedger) Balance(addr string) uint64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.balances[addr]
}

// Burn removes tokens from an address.
func (l *DAOTokenLedger) Burn(addr string, amount uint64) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	bal := l.balances[addr]
	if bal < amount {
		return errors.New("insufficient balance")
	}
	l.balances[addr] = bal - amount
	return nil
}
