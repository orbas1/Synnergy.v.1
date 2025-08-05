package core

import (
	"sync"
	nodes "synnergy/internal/nodes"
)

// ForensicNode captures lightweight transaction data and network traces for
// later inspection. It stores entries in memory; a production implementation
// would persist this information to durable storage.
type ForensicNode struct {
	mu     sync.RWMutex
	txs    []nodes.TransactionLite
	traces []nodes.NetworkTrace
}

// Ensure ForensicNode implements nodes.ForensicNodeInterface.
var _ nodes.ForensicNodeInterface = (*ForensicNode)(nil)

// NewForensicNode creates a ForensicNode with empty buffers.
func NewForensicNode() *ForensicNode {
	return &ForensicNode{}
}

// RecordTransaction stores a minimal representation of a transaction.
func (f *ForensicNode) RecordTransaction(tx nodes.TransactionLite) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.txs = append(f.txs, tx)
	return nil
}

// RecordNetworkTrace stores a network level event for later analysis.
func (f *ForensicNode) RecordNetworkTrace(trace nodes.NetworkTrace) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.traces = append(f.traces, trace)
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
