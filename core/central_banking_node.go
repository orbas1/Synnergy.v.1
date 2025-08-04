package core

// CentralBankingNode models a node operated by a central bank.
type CentralBankingNode struct {
	*Node
	MonetaryPolicy string
}

// NewCentralBankingNode creates a central banking node with the given policy.
func NewCentralBankingNode(id, addr string, ledger *Ledger, policy string) *CentralBankingNode {
	return &CentralBankingNode{
		Node:           NewNode(id, addr, ledger),
		MonetaryPolicy: policy,
	}
}

// UpdatePolicy updates the node's monetary policy description.
func (n *CentralBankingNode) UpdatePolicy(policy string) {
	n.MonetaryPolicy = policy
}

// Mint increases the treasury within the attached ledger.
func (n *CentralBankingNode) Mint(to string, amount uint64) {
	n.Ledger.Credit(to, amount)
}
