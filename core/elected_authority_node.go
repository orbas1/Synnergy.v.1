package core

import (
	"sync"
	"time"
)

// ElectedAuthorityNode represents an authority node with a fixed term.
// A mutex guards the term so concurrent renewals are safe when multiple
// governance operations attempt to extend a node's tenure.
type ElectedAuthorityNode struct {
	*AuthorityNode
	mu      sync.RWMutex
	TermEnd time.Time
}

// NewElectedAuthorityNode creates a new elected authority node with the given
// term duration and an empty vote set.
func NewElectedAuthorityNode(addr, role string, term time.Duration) *ElectedAuthorityNode {
	node := &AuthorityNode{Address: addr, Role: role, Votes: make(map[string]bool)}
	return &ElectedAuthorityNode{AuthorityNode: node, TermEnd: time.Now().Add(term)}
}

// IsActive returns true if the node's term has not expired.
func (n *ElectedAuthorityNode) IsActive(now time.Time) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return now.Before(n.TermEnd)
}

// RenewTerm extends the node's term. The requester must be a DAO admin.
// This ensures that only authorized governance participants can prolong an
// authority node's tenure.
func (n *ElectedAuthorityNode) RenewTerm(requester string, dao *DAO, add time.Duration) error {
	if dao == nil || !dao.IsAdmin(requester) {
		return errUnauthorized
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.TermEnd = n.TermEnd.Add(add)
	return nil
}
