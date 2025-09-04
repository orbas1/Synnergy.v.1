package core

import (
	"fmt"

	"synnergy/internal/tokens"
)

// CentralBankingNode models a node operated by a central bank. It may manage
// the supply of a CBDC token but cannot mint the fixed-supply native SYN coin.
type CentralBankingNode struct {
	*Node
	MonetaryPolicy string
	CBDCToken      *tokens.SYN10Token
}

// NewCentralBankingNode creates a central banking node with the given policy
// and associated CBDC token.
func NewCentralBankingNode(id, addr string, ledger *Ledger, policy string, token *tokens.SYN10Token) *CentralBankingNode {
	return &CentralBankingNode{
		Node:           NewNode(id, addr, ledger),
		MonetaryPolicy: policy,
		CBDCToken:      token,
	}
}

// UpdatePolicy updates the node's monetary policy description.
func (n *CentralBankingNode) UpdatePolicy(policy string) { n.MonetaryPolicy = policy }

// MintCBDC mints new units of the CBDC token. It rejects zero amounts. The
// native SYN coin supply remains unaffected and permanently capped.
func (n *CentralBankingNode) MintCBDC(to string, amount uint64) error {
	if amount == 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	return n.CBDCToken.Mint(to, amount)
}

// Mint is retained for backwards compatibility but now forwards to MintCBDC.
// This preserves existing opcode mappings while ensuring SYN coin supply is
// never altered through this method.
func (n *CentralBankingNode) Mint(to string, amount uint64) error { return n.MintCBDC(to, amount) }
