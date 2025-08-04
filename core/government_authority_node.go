package core

// GovernmentAuthorityNode represents a regulator-operated authority node.
type GovernmentAuthorityNode struct {
	*AuthorityNode
	Department string
}

// NewGovernmentAuthorityNode creates a new government authority node.
func NewGovernmentAuthorityNode(addr, role, department string) *GovernmentAuthorityNode {
	node := &AuthorityNode{Address: addr, Role: role, Votes: make(map[string]bool)}
	return &GovernmentAuthorityNode{AuthorityNode: node, Department: department}
}
