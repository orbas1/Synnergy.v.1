package core

// BankInstitutionalNode represents a banking institution participating in the network.
type BankInstitutionalNode struct {
	*Node
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
	n.Institutions[name] = true
}

// IsRegistered checks if an institution is registered.
func (n *BankInstitutionalNode) IsRegistered(name string) bool {
	return n.Institutions[name]
}
