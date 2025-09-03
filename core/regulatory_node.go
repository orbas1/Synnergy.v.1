package core

import (
	"strings"
	"sync"
)

// RegulatoryNode represents a regulator-operated node overseeing transactions.
type RegulatoryNode struct {
	ID      string
	Manager *RegulatoryManager
	mu      sync.RWMutex
	logs    map[string][]string
}

// NewRegulatoryNode creates a new RegulatoryNode.
func NewRegulatoryNode(id string, mgr *RegulatoryManager) *RegulatoryNode {
	return &RegulatoryNode{
		ID:      id,
		Manager: mgr,
		logs:    make(map[string][]string),
	}
}

// ApproveTransaction checks a transaction against registered regulations.
func (n *RegulatoryNode) ApproveTransaction(tx Transaction) bool {
	violations := n.Manager.EvaluateTransaction(tx)
	if len(violations) > 0 {
		n.FlagEntity(tx.From, strings.Join(violations, ", "))
		return false
	}
	return true
}

// FlagEntity records a regulatory flag for an address.
func (n *RegulatoryNode) FlagEntity(addr, reason string) {
	n.mu.Lock()
	n.logs[addr] = append(n.logs[addr], reason)
	n.mu.Unlock()
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
