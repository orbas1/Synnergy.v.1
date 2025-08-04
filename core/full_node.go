package core

import "synnergy/nodes"

// FullNodeMode specifies the storage strategy of a full node.
type FullNodeMode int

const (
	// FullNodeModeArchive stores the entire blockchain history.
	FullNodeModeArchive FullNodeMode = iota
	// FullNodeModePruned retains only recent blocks to save space.
	FullNodeModePruned
)

// FullNode represents a standard validating node storing the full chain.
type FullNode struct {
	*BaseNode
	Mode FullNodeMode
}

// NewFullNode creates a full node with the given mode.
func NewFullNode(id nodes.Address, mode FullNodeMode) *FullNode {
	return &FullNode{BaseNode: NewBaseNode(id), Mode: mode}
}
