package banknodes

import (
	"crypto/ed25519"
	"synnergy/internal/nodes"
)

// BankInstitutionalNode defines behaviour for institutional banking nodes.
// It extends the generic NodeInterface with methods to manage participating institutions.
type BankInstitutionalNode interface {
	nodes.NodeInterface
	// RegisterInstitution registers a new institution using signature verification.
	RegisterInstitution(addr, name string, sig []byte, pubKey ed25519.PublicKey) error
	// RemoveInstitution removes an institution via a signed request.
	RemoveInstitution(addr, name string, sig []byte, pubKey ed25519.PublicKey) error
	// ListInstitutions returns all currently registered institutions.
	ListInstitutions() []string
	// IsRegistered checks whether an institution is already registered.
	IsRegistered(name string) bool
}

// CentralBankingNode defines operations for nodes operated by central banks.
type CentralBankingNode interface {
	nodes.NodeInterface
	// UpdatePolicy updates the node's monetary policy guidance.
	UpdatePolicy(policy string)
	// MintCBDC creates new CBDC tokens for the target account.
	MintCBDC(to string, amount uint64) error
}

// CustodialNode models a node that holds assets on behalf of users.
type CustodialNode interface {
	nodes.NodeInterface
	// Custody records assets held for a specific user.
	Custody(user string, amount uint64)
	// Release transfers assets back to the user if sufficient.
	Release(user string, amount uint64) error
}
