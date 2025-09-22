package nodes

import (
	"errors"
	"sync"
	"time"
)

// TransactionLite represents the minimal transaction information captured by
// a forensic node without requiring full ledger access.
type TransactionLite struct {
	Hash      string
	From      string
	To        string
	Value     uint64
	Timestamp time.Time
}

// NetworkTrace captures a single network level event for later forensic
// analysis.
type NetworkTrace struct {
	PeerID    string
	Event     string
	Timestamp time.Time
}

// ForensicNodeInterface defines behaviour for nodes that record lightweight
// transaction details and network traces for offline analysis.
type ForensicNodeInterface interface {
	// RecordTransaction stores a minimal representation of a transaction for
	// later inspection.
	RecordTransaction(tx TransactionLite) error
	// RecordNetworkTrace stores a network trace or event for later analysis.
	RecordNetworkTrace(trace NetworkTrace) error
	// Transactions returns the list of recorded transactions.
	Transactions() []TransactionLite
	// NetworkTraces returns the list of captured network traces.
	NetworkTraces() []NetworkTrace
}

// ForensicNode provides an in-memory implementation of ForensicNodeInterface
// backed by a BasicNode for standard lifecycle management. Data recorded by the
// node is stored in-memory and therefore only intended for testing and
// development networks.
type ForensicNode struct {
	*BasicNode
	mu              sync.RWMutex
	txs             []TransactionLite
	traces          []NetworkTrace
	maxTransactions int
	maxTraces       int
	lastTransaction time.Time
	lastTrace       time.Time
}

// ForensicStats provides operational insight for monitoring and CLI tooling.
type ForensicStats struct {
	TransactionCount int
	TraceCount       int
	LastTransaction  time.Time
	LastTrace        time.Time
}

const (
	defaultForensicLimit = 1024
)

type forensicConfig struct {
	maxTransactions int
	maxTraces       int
}

// ForensicOption configures optional behaviour for a forensic node.
type ForensicOption func(*forensicConfig)

// WithMaxTransactionRecords sets the maximum number of transaction snapshots
// the node will retain in memory. A non-positive value reverts to the default.
func WithMaxTransactionRecords(limit int) ForensicOption {
	return func(cfg *forensicConfig) {
		if limit > 0 {
			cfg.maxTransactions = limit
		}
	}
}

// WithMaxTraceRecords sets the maximum number of network traces retained.
// A non-positive value reverts to the default.
func WithMaxTraceRecords(limit int) ForensicOption {
	return func(cfg *forensicConfig) {
		if limit > 0 {
			cfg.maxTraces = limit
		}
	}
}

// NewForensicNode creates a new forensic node with the provided identifier and
// optional configuration.
func NewForensicNode(id Address, opts ...ForensicOption) *ForensicNode {
	cfg := forensicConfig{
		maxTransactions: defaultForensicLimit,
		maxTraces:       defaultForensicLimit,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &ForensicNode{
		BasicNode:       NewBasicNode(id),
		maxTransactions: cfg.maxTransactions,
		maxTraces:       cfg.maxTraces,
	}
}

// RecordTransaction saves a transaction snapshot for later inspection.
func (n *ForensicNode) RecordTransaction(tx TransactionLite) error {
	if tx.Hash == "" {
		return errors.New("transaction hash is required")
	}
	if tx.Timestamp.IsZero() {
		tx.Timestamp = time.Now().UTC()
	}

	n.mu.Lock()
	n.txs = append(n.txs, tx)
	if n.maxTransactions > 0 && len(n.txs) > n.maxTransactions {
		n.txs = append([]TransactionLite(nil), n.txs[len(n.txs)-n.maxTransactions:]...)
	}
	n.lastTransaction = tx.Timestamp
	n.mu.Unlock()
	return nil
}

// RecordNetworkTrace stores a network trace event.
func (n *ForensicNode) RecordNetworkTrace(trace NetworkTrace) error {
	if trace.PeerID == "" {
		return errors.New("peer id is required")
	}
	if trace.Timestamp.IsZero() {
		trace.Timestamp = time.Now().UTC()
	}

	n.mu.Lock()
	n.traces = append(n.traces, trace)
	if n.maxTraces > 0 && len(n.traces) > n.maxTraces {
		n.traces = append([]NetworkTrace(nil), n.traces[len(n.traces)-n.maxTraces:]...)
	}
	n.lastTrace = trace.Timestamp
	n.mu.Unlock()
	return nil
}

// Transactions returns a copy of all recorded transactions.
func (n *ForensicNode) Transactions() []TransactionLite {
	n.mu.RLock()
	out := make([]TransactionLite, len(n.txs))
	copy(out, n.txs)
	n.mu.RUnlock()
	return out
}

// NetworkTraces returns a copy of all captured network traces.
func (n *ForensicNode) NetworkTraces() []NetworkTrace {
	n.mu.RLock()
	out := make([]NetworkTrace, len(n.traces))
	copy(out, n.traces)
	n.mu.RUnlock()
	return out
}

// Stats returns operational counters for monitoring or CLI display.
func (n *ForensicNode) Stats() ForensicStats {
	n.mu.RLock()
	stats := ForensicStats{
		TransactionCount: len(n.txs),
		TraceCount:       len(n.traces),
		LastTransaction:  n.lastTransaction,
		LastTrace:        n.lastTrace,
	}
	n.mu.RUnlock()
	return stats
}

// Ensure ForensicNode implements ForensicNodeInterface.
var _ ForensicNodeInterface = (*ForensicNode)(nil)
