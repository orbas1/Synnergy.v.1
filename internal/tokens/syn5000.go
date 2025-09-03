package tokens

import (
	"errors"
	"sync"
)

// SYN5000Token records balances across multiple chains and supports
// cross-chain transfers.
type SYN5000Token struct {
	mu       sync.RWMutex
	balances map[string]map[string]uint64 // chain -> address -> balance
}

// NewSYN5000Token creates an empty multi-chain token registry.
func NewSYN5000Token() *SYN5000Token {
	return &SYN5000Token{balances: make(map[string]map[string]uint64)}
}

// Mint credits an address on a specific chain.
func (t *SYN5000Token) Mint(chain, addr string, amount uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.balances[chain] == nil {
		t.balances[chain] = make(map[string]uint64)
	}
	t.balances[chain][addr] += amount
}

// Transfer moves funds from one chain/address to another.
func (t *SYN5000Token) Transfer(fromChain, fromAddr, toChain, toAddr string, amount uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	fromBal := t.balances[fromChain][fromAddr]
	if fromBal < amount {
		return errors.New("insufficient balance")
	}
	t.balances[fromChain][fromAddr] = fromBal - amount
	if t.balances[toChain] == nil {
		t.balances[toChain] = make(map[string]uint64)
	}
	t.balances[toChain][toAddr] += amount
	return nil
}

// Balance returns the balance for an address on a chain.
func (t *SYN5000Token) Balance(chain, addr string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.balances[chain] == nil {
		return 0
	}
	return t.balances[chain][addr]
}
