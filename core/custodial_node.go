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
	Relayers map[string]struct{}
}

// NewCustodialNode creates a custodial node instance.
func NewCustodialNode(id, addr string, ledger *Ledger) *CustodialNode {
	return &CustodialNode{
		Node:     NewNode(id, addr, ledger),
		Holdings: make(map[string]uint64),
		Relayers: make(map[string]struct{}),
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
func (n *CustodialNode) AuthorizeRelayer(addr string) {
	n.mu.Lock()
	n.Relayers[addr] = struct{}{}
	n.mu.Unlock()
}

// RevokeRelayer removes an address from the whitelist.
func (n *CustodialNode) RevokeRelayer(addr string) {
	n.mu.Lock()
	delete(n.Relayers, addr)
	n.mu.Unlock()
}

// IsRelayerAuthorized returns true if the address is whitelisted.
func (n *CustodialNode) IsRelayerAuthorized(addr string) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	_, ok := n.Relayers[addr]
	return ok
}

// Release transfers assets back to a user if sufficient and credits the
// underlying ledger. Only authorized relayers can trigger a release. It
// returns an error on insufficient holdings or unauthorized relayers.
func (n *CustodialNode) Release(user string, amount uint64, relayer string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if _, ok := n.Relayers[relayer]; !ok {
		return errors.New("relayer not authorized")
	}
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
