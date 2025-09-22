package synnergy

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// RegulatoryNode represents a regulator-operated node overseeing transactions.
type RegulatoryNode struct {
	ID        string
	Manager   *RegulatoryManager
	mu        sync.RWMutex
	logs      map[string][]string
	approvals map[string]time.Time
}

// NewRegulatoryNode creates a new RegulatoryNode.
func NewRegulatoryNode(id string, mgr *RegulatoryManager) *RegulatoryNode {
	return &RegulatoryNode{
		ID:        id,
		Manager:   mgr,
		logs:      make(map[string][]string),
		approvals: make(map[string]time.Time),
	}
}

// ApproveTransaction checks a transaction against registered regulations and
// records compliance decisions for later audits.
func (n *RegulatoryNode) ApproveTransaction(tx Transaction) bool {
	if n.Manager == nil {
		n.FlagEntity(tx.From, "no regulatory manager configured")
		return false
	}
	res := n.Manager.EvaluateTransactionDetailed(tx)
	if len(res.Violations) > 0 {
		reasons := make([]string, 0, len(res.Violations))
		for _, v := range res.Violations {
			reasons = append(reasons, fmt.Sprintf("%s: %s", v.RegulationID, v.Reason))
		}
		n.FlagEntity(tx.From, strings.Join(reasons, "; "))
		return false
	}
	n.recordApproval(tx.From)
	return true
}

func (n *RegulatoryNode) recordApproval(addr string) {
	n.mu.Lock()
	n.approvals[addr] = time.Now().UTC()
	n.mu.Unlock()
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

// LastApproval returns the timestamp of the most recent approval for an address.
func (n *RegulatoryNode) LastApproval(addr string) (time.Time, bool) {
	n.mu.RLock()
	ts, ok := n.approvals[addr]
	n.mu.RUnlock()
	return ts, ok
}
