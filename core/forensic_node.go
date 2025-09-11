package core

import (
	"sync"

	nodes "synnergy/internal/nodes"
)

// ForensicNode captures lightweight transaction data and network traces for
// later inspection. It stores a bounded number of entries in memory so a noisy
// peer cannot exhaust resources. Older records are pruned in FIFO order once
// the capacity is exceeded.
type ForensicNode struct {
	mu     sync.RWMutex
	txs    []nodes.TransactionLite
	traces []nodes.NetworkTrace
	limit  int
}

// Ensure ForensicNode implements nodes.ForensicNodeInterface.
var _ nodes.ForensicNodeInterface = (*ForensicNode)(nil)

const DefaultForensicLimit = 1000

// NewForensicNode creates a ForensicNode with empty buffers and the default
// capacity limit.
func NewForensicNode() *ForensicNode {
	return &ForensicNode{limit: DefaultForensicLimit}
}

// NewForensicNodeWithLimit creates a ForensicNode that retains at most limit
// entries for each log. If limit <= 0 the default is used.
func NewForensicNodeWithLimit(limit int) *ForensicNode {
	if limit <= 0 {
		limit = DefaultForensicLimit
	}
	return &ForensicNode{limit: limit}
}

// RecordTransaction stores a minimal representation of a transaction.
func (f *ForensicNode) RecordTransaction(tx nodes.TransactionLite) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.txs = append(f.txs, tx)
	if len(f.txs) > f.limit {
		f.txs = append([]nodes.TransactionLite(nil), f.txs[len(f.txs)-f.limit:]...)
	}
	return nil
}

// RecordNetworkTrace stores a network level event for later analysis.
func (f *ForensicNode) RecordNetworkTrace(trace nodes.NetworkTrace) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.traces = append(f.traces, trace)
	if len(f.traces) > f.limit {
		f.traces = append([]nodes.NetworkTrace(nil), f.traces[len(f.traces)-f.limit:]...)
	}
	return nil
}

// Transactions returns a copy of recorded transactions.
func (f *ForensicNode) Transactions() []nodes.TransactionLite {
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := make([]nodes.TransactionLite, len(f.txs))
	copy(out, f.txs)
	return out
}

// NetworkTraces returns a copy of captured network traces.
func (f *ForensicNode) NetworkTraces() []nodes.NetworkTrace {
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := make([]nodes.NetworkTrace, len(f.traces))
	copy(out, f.traces)
	return out
}
