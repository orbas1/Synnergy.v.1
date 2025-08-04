package core

import (
	"errors"
	"sync"
)

// Ledger maintains account balances.
type Ledger struct {
	mu       sync.RWMutex
	balances map[string]uint64
}

// NewLedger creates a new empty ledger.
func NewLedger() *Ledger {
	return &Ledger{balances: make(map[string]uint64)}
}

// GetBalance returns the balance for a given address.
func (l *Ledger) GetBalance(addr string) uint64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.balances[addr]
}

// Credit adds funds to an address.
func (l *Ledger) Credit(addr string, amount uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.balances[addr] += amount
}

// ApplyTransaction applies a transaction to the ledger, deducting both amount
// and fee from the sender. It returns an error if the sender lacks sufficient
// funds.
func (l *Ledger) ApplyTransaction(tx *Transaction) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	total := tx.Amount + tx.Fee
	if l.balances[tx.From] < total {
		return errors.New("insufficient funds")
	}
	l.balances[tx.From] -= total
	l.balances[tx.To] += tx.Amount
	return nil
}
