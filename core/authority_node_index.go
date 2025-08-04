package core

// AuthorityNodeIndex maintains a lookup of authority nodes by address.
type AuthorityNodeIndex struct {
	nodes map[string]*AuthorityNode
}

// NewAuthorityNodeIndex returns an initialised AuthorityNodeIndex.
func NewAuthorityNodeIndex() *AuthorityNodeIndex {
	return &AuthorityNodeIndex{nodes: make(map[string]*AuthorityNode)}
}

// Add inserts or replaces an authority node in the index.
func (idx *AuthorityNodeIndex) Add(node *AuthorityNode) {
	if idx.nodes == nil {
		idx.nodes = make(map[string]*AuthorityNode)
	}
	idx.nodes[node.Address] = node
}

// Get retrieves an authority node by address.
func (idx *AuthorityNodeIndex) Get(addr string) (*AuthorityNode, bool) {
	n, ok := idx.nodes[addr]
	return n, ok
}

// Remove deletes an authority node from the index by address.
func (idx *AuthorityNodeIndex) Remove(addr string) {
	delete(idx.nodes, addr)
}

// List returns all authority nodes in the index.
func (idx *AuthorityNodeIndex) List() []*AuthorityNode {
	out := make([]*AuthorityNode, 0, len(idx.nodes))
	for _, n := range idx.nodes {
		out = append(out, n)
	}
	return out
}
