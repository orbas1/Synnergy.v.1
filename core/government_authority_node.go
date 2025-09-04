package core

import "errors"

// GovernmentAuthorityNode represents a regulator-operated authority node.
// It explicitly lacks capabilities to mint SYN coins or modify monetary policy.
type GovernmentAuthorityNode struct {
	*AuthorityNode
	Department string
}

// NewGovernmentAuthorityNode creates a new government authority node.
func NewGovernmentAuthorityNode(addr, role, department string) *GovernmentAuthorityNode {
	node := &AuthorityNode{Address: addr, Role: role, Votes: make(map[string]bool)}
	return &GovernmentAuthorityNode{AuthorityNode: node, Department: department}
}

// MintSYN always returns an error as government nodes cannot mint the native
// SYN coin which has a fixed supply.
func (n *GovernmentAuthorityNode) MintSYN(to string, amount uint64) error {
	return errors.New("government authority nodes cannot mint SYN coins")
}

// UpdateMonetaryPolicy always returns an error as government nodes cannot change
// monetary policy. Only central bank nodes may perform such actions.
func (n *GovernmentAuthorityNode) UpdateMonetaryPolicy(policy string) error {
	return errors.New("government authority nodes cannot modify monetary policy")
}
