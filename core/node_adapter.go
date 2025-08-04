package core

import "synnergy/nodes"

// NodeAdapter adapts the consensus Node to the generic nodes.NodeInterface.
type NodeAdapter struct {
	node *Node
	*BaseNode
}

// NewNodeAdapter wraps the provided Node with a BaseNode implementing nodes.NodeInterface.
func NewNodeAdapter(n *Node) *NodeAdapter {
	return &NodeAdapter{
		node:     n,
		BaseNode: NewBaseNode(nodes.Address(n.ID)),
	}
}

var _ nodes.NodeInterface = (*NodeAdapter)(nil)
