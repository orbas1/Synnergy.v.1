package core

import "synnergy/nodes"

// LightNode maintains a minimal view of the chain using block headers only.
type LightNode struct {
	*BaseNode
	headers []nodes.BlockHeader
}

// NewLightNode constructs a light node with no headers.
func NewLightNode(id nodes.Address) *LightNode {
	return &LightNode{BaseNode: NewBaseNode(id)}
}

// AddHeader appends a header to the local view.
func (n *LightNode) AddHeader(h nodes.BlockHeader) { n.headers = append(n.headers, h) }

// LatestHeader returns the most recently added header.
func (n *LightNode) LatestHeader() (nodes.BlockHeader, bool) {
	if len(n.headers) == 0 {
		return nodes.BlockHeader{}, false
	}
	return n.headers[len(n.headers)-1], true
}

// Headers returns a copy of all stored block headers.
func (n *LightNode) Headers() []nodes.BlockHeader {
	out := make([]nodes.BlockHeader, len(n.headers))
	copy(out, n.headers)
	return out
}
