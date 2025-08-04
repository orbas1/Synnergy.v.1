package nodes

import (
	"time"

	an "synnergy/nodes/authority_nodes"
)

// ElectedAuthorityNode represents an authority node with a fixed term.
type ElectedAuthorityNode struct {
	*an.AuthorityNode
	TermEnd time.Time
}

// NewElectedAuthorityNode creates a new elected authority node with the given term duration.
func NewElectedAuthorityNode(addr Address, role string, term time.Duration) *ElectedAuthorityNode {
	node := &an.AuthorityNode{Address: string(addr), Role: role, Votes: make(map[string]bool)}
	return &ElectedAuthorityNode{AuthorityNode: node, TermEnd: time.Now().Add(term)}
}

// IsActive returns true if the node's term has not expired.
func (n *ElectedAuthorityNode) IsActive(now time.Time) bool {
	return now.Before(n.TermEnd)
}
