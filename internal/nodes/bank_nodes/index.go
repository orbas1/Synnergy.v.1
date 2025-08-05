package banknodes

import "synnergy/internal/nodes"

// BankInstitutionalNode defines behaviour for institutional banking nodes.
// It extends the generic NodeInterface with methods to manage participating institutions.
type BankInstitutionalNode interface {
	nodes.NodeInterface
	// RegisterInstitution registers a new institution by name.
	RegisterInstitution(name string)
	// IsRegistered checks whether an institution is already registered.
	IsRegistered(name string) bool
}

// CentralBankingNode defines operations for nodes operated by central banks.
type CentralBankingNode interface {
	nodes.NodeInterface
	// UpdatePolicy updates the node's monetary policy guidance.
	UpdatePolicy(policy string)
	// Mint credits the given amount to the target account within the ledger.
	Mint(to string, amount uint64)
}

// CustodialNode models a node that holds assets on behalf of users.
type CustodialNode interface {
	nodes.NodeInterface
	// Custody records assets held for a specific user.
	Custody(user string, amount uint64)
	// Release transfers assets back to the user if sufficient and returns success.
	Release(user string, amount uint64) bool
}
