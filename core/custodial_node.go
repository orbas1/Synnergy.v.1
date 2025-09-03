package core

import (
	"errors"
	"sync"
)

// CustodialNode holds assets on behalf of users. Operations are guarded by a
// mutex so they are safe for concurrent use by wallet and node processes.
type CustodialNode struct {
	*Node
	mu       sync.RWMutex
	Holdings map[string]uint64
}

// NewCustodialNode creates a custodial node instance.
func NewCustodialNode(id, addr string, ledger *Ledger) *CustodialNode {
	return &CustodialNode{
		Node:     NewNode(id, addr, ledger),
		Holdings: make(map[string]uint64),
	}
}

// Custody records assets held for a user.
func (n *CustodialNode) Custody(user string, amount uint64) {
	n.mu.Lock()
	n.Holdings[user] += amount
	n.mu.Unlock()
}

// Release transfers assets back to a user if sufficient and credits the
// underlying ledger. It returns an error on insufficient holdings.
func (n *CustodialNode) Release(user string, amount uint64) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.Holdings[user] < amount {
		return errors.New("insufficient holdings")
	}
	n.Holdings[user] -= amount
	n.Ledger.Credit(user, amount)
	return nil
}

// Balance returns the holdings for a given user.
func (n *CustodialNode) Balance(user string) uint64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Holdings[user]
}
