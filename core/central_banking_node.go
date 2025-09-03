package core

import "fmt"

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

// Mint increases the treasury within the attached ledger. It rejects mints that
// would exceed the remaining supply defined in coin.go.
func (n *CentralBankingNode) Mint(to string, amount uint64) error {
	if amount == 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	h, _ := n.Ledger.Head()
	if RemainingSupply(uint64(h)) < amount {
		return fmt.Errorf("mint exceeds remaining supply")
	}
	n.Ledger.Credit(to, amount)
	return nil
}
