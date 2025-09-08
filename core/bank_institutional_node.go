package core

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"
)

// BankInstitutionalNode represents a banking institution participating in the network.
type BankInstitutionalNode struct {
	*Node
	mu           sync.RWMutex
	Institutions map[string]bool
}

// NewBankInstitutionalNode creates a new institutional banking node.
func NewBankInstitutionalNode(id, addr string, ledger *Ledger) *BankInstitutionalNode {
	return &BankInstitutionalNode{
		Node:         NewNode(id, addr, ledger),
		Institutions: make(map[string]bool),
	}
}

// RegisterInstitution records an institution after verifying the signature.
// addr must match pubKey and the signature must cover "register:name".
func (n *BankInstitutionalNode) RegisterInstitution(addr, name string, sig []byte, pubKey ed25519.PublicKey) error {
	if hex.EncodeToString(pubKey) != addr {
		return errors.New("address mismatch")
	}
	msg := "register:" + name
	if !ed25519.Verify(pubKey, []byte(msg), sig) {
		return errors.New("invalid signature")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.Institutions == nil {
		n.Institutions = make(map[string]bool)
	}
	n.Institutions[name] = true
	return nil
}

// RemoveInstitution deletes an institution after verifying the signature.
// addr must match pubKey and the signature must cover "remove:name".
func (n *BankInstitutionalNode) RemoveInstitution(addr, name string, sig []byte, pubKey ed25519.PublicKey) error {
	if hex.EncodeToString(pubKey) != addr {
		return errors.New("address mismatch")
	}
	msg := "remove:" + name
	if !ed25519.Verify(pubKey, []byte(msg), sig) {
		return errors.New("invalid signature")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.Institutions, name)
	return nil
}

// ListInstitutions returns a snapshot of registered institutions.
func (n *BankInstitutionalNode) ListInstitutions() []string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	out := make([]string, 0, len(n.Institutions))
	for name := range n.Institutions {
		out = append(out, name)
	}
	return out
}

// IsRegistered checks if an institution is registered.
func (n *BankInstitutionalNode) IsRegistered(name string) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Institutions[name]
}

// MarshalJSON returns a serialisable representation of the node.
func (n *BankInstitutionalNode) MarshalJSON() ([]byte, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	type alias BankInstitutionalNode
	return json.Marshal(&struct {
		Institutions []string `json:"institutions"`
		*alias
	}{
		Institutions: n.ListInstitutions(),
		alias:        (*alias)(n),
	})
}
