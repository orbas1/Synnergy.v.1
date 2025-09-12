package core

import (
	"crypto/ecdsa"
	"errors"
	"sync"
)

// RegulatoryNode represents a regulator-operated node overseeing transactions.
type RegulatoryNode struct {
	ID      string
	Manager *RegulatoryManager
	mu      sync.RWMutex
	logs    map[string][]string
	wallets map[string]*ecdsa.PublicKey
}

// NewRegulatoryNode creates a new RegulatoryNode.
func NewRegulatoryNode(id string, mgr *RegulatoryManager) *RegulatoryNode {
	return &RegulatoryNode{
		ID:      id,
		Manager: mgr,
		logs:    make(map[string][]string),
		wallets: make(map[string]*ecdsa.PublicKey),
	}
}

// ApproveTransaction checks a transaction against registered regulations and
// returns an error if any rules are violated.
func (n *RegulatoryNode) ApproveTransaction(tx Transaction) error {
	if n.Manager == nil {
		return nil
	}
	pub, ok := n.wallets[tx.From]
	if !ok {
		n.FlagEntity(tx.From, "unknown wallet")
		return errors.New("unknown wallet")
	}
	if !tx.Verify(pub) {
		n.FlagEntity(tx.From, "invalid signature")
		return errors.New("invalid signature")
	}
	if err := n.Manager.ValidateTransaction(tx); err != nil {
		n.FlagEntity(tx.From, err.Error())
		return err
	}
	return nil
}

// FlagEntity records a regulatory flag for an address.
func (n *RegulatoryNode) FlagEntity(addr, reason string) error {
	if reason == "" {
		return errors.New("reason required")
	}
	n.mu.Lock()
	n.logs[addr] = append(n.logs[addr], reason)
	n.mu.Unlock()
	return nil
}

// Logs returns all flags recorded for an address.
func (n *RegulatoryNode) Logs(addr string) []string {
	n.mu.RLock()
	entries := n.logs[addr]
	out := make([]string, len(entries))
	copy(out, entries)
	n.mu.RUnlock()
	return out
}

// ClearLogs removes all flags recorded for an address.
func (n *RegulatoryNode) ClearLogs(addr string) {
	n.mu.Lock()
	delete(n.logs, addr)
	n.mu.Unlock()
}

// RegisterWallet associates a wallet's public key with its address so that
// future transactions can be signature verified.
func (n *RegulatoryNode) RegisterWallet(w *Wallet) {
	n.mu.Lock()
	pub := w.PrivateKey.PublicKey
	n.wallets[w.Address] = &pub
	n.mu.Unlock()
}
