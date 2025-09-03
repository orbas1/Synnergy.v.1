package core

import (
        "encoding/json"
        "sync"
)

// AuthorityNodeIndex maintains a lookup of authority nodes by address.
type AuthorityNodeIndex struct {
	mu    sync.RWMutex
	nodes map[string]*AuthorityNode
}

// NewAuthorityNodeIndex returns an initialised AuthorityNodeIndex.
func NewAuthorityNodeIndex() *AuthorityNodeIndex {
	return &AuthorityNodeIndex{nodes: make(map[string]*AuthorityNode)}
}

// Add inserts or replaces an authority node in the index.
func (idx *AuthorityNodeIndex) Add(node *AuthorityNode) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	if idx.nodes == nil {
		idx.nodes = make(map[string]*AuthorityNode)
	}
	idx.nodes[node.Address] = node
}

// Get retrieves an authority node by address.
func (idx *AuthorityNodeIndex) Get(addr string) (*AuthorityNode, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	n, ok := idx.nodes[addr]
	return n, ok
}

// Remove deletes an authority node from the index by address.
func (idx *AuthorityNodeIndex) Remove(addr string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	delete(idx.nodes, addr)
}

// List returns all authority nodes in the index.
func (idx *AuthorityNodeIndex) List() []*AuthorityNode {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	out := make([]*AuthorityNode, 0, len(idx.nodes))
	for _, n := range idx.nodes {
		out = append(out, n)
	}
	return out
}

// Snapshot returns a copy of the underlying map for safe iteration.
func (idx *AuthorityNodeIndex) Snapshot() map[string]*AuthorityNode {
        idx.mu.RLock()
        defer idx.mu.RUnlock()
        out := make(map[string]*AuthorityNode, len(idx.nodes))
        for k, v := range idx.nodes {
                out[k] = v
        }
        return out
}

// MarshalJSON emits the index as an array for easier consumption via the CLI.
func (idx *AuthorityNodeIndex) MarshalJSON() ([]byte, error) {
        return json.Marshal(idx.List())
}
