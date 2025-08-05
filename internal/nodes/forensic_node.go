package nodes

import "time"

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
