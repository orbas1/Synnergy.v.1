package core

import "sync"

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

// ApplyTransaction applies a transaction to the ledger.
func (l *Ledger) ApplyTransaction(tx *Transaction) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.balances[tx.From] -= tx.Amount
	l.balances[tx.To] += tx.Amount
}
