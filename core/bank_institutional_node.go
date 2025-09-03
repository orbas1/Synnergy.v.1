package core

import (
        "encoding/json"
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

// RegisterInstitution registers a participating institution by name.
func (n *BankInstitutionalNode) RegisterInstitution(name string) {
        n.mu.Lock()
        defer n.mu.Unlock()
        if n.Institutions == nil {
                n.Institutions = make(map[string]bool)
        }
        n.Institutions[name] = true
}

// RemoveInstitution removes a participating institution.
func (n *BankInstitutionalNode) RemoveInstitution(name string) {
        n.mu.Lock()
        defer n.mu.Unlock()
        delete(n.Institutions, name)
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
        return json.Marshal(&struct{
                Institutions []string `json:"institutions"`
                *alias
        }{
                Institutions: n.ListInstitutions(),
                alias:        (*alias)(n),
        })
}
