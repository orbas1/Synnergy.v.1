package authority_nodes

// AuthorityNode represents a node eligible for governance actions.
type AuthorityNode struct {
	Address string
	Role    string
	Votes   map[string]bool // voter address -> approved
}

// Index maintains a lookup of authority nodes by address.
type Index struct {
	nodes map[string]*AuthorityNode
}

// NewIndex returns an initialised Index.
func NewIndex() *Index {
	return &Index{nodes: make(map[string]*AuthorityNode)}
}

// Add inserts or replaces an authority node in the index.
func (idx *Index) Add(node *AuthorityNode) {
	if idx.nodes == nil {
		idx.nodes = make(map[string]*AuthorityNode)
	}
	idx.nodes[node.Address] = node
}

// Get retrieves an authority node by address.
func (idx *Index) Get(addr string) (*AuthorityNode, bool) {
	n, ok := idx.nodes[addr]
	return n, ok
}

// Remove deletes an authority node from the index by address.
func (idx *Index) Remove(addr string) {
	delete(idx.nodes, addr)
}

// List returns all authority nodes in the index.
func (idx *Index) List() []*AuthorityNode {
	out := make([]*AuthorityNode, 0, len(idx.nodes))
	for _, n := range idx.nodes {
		out = append(out, n)
	}
	return out
}
