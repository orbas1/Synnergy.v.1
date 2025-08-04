package core

// CustodialNode holds assets on behalf of users.
type CustodialNode struct {
	*Node
	Holdings map[string]uint64
}

// NewCustodialNode creates a custodial node instance.
func NewCustodialNode(id, addr string, ledger *Ledger) *CustodialNode {
	return &CustodialNode{
		Node:     NewNode(id, addr, ledger),
		Holdings: make(map[string]uint64),
	}
}

// Custody records assets held for a user.
func (n *CustodialNode) Custody(user string, amount uint64) {
	n.Holdings[user] += amount
}

// Release transfers assets back to a user if sufficient.
func (n *CustodialNode) Release(user string, amount uint64) bool {
	if n.Holdings[user] < amount {
		return false
	}
	n.Holdings[user] -= amount
	n.Ledger.Credit(user, amount)
	return true
}
