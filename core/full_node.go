package core

import "synnergy/internal/nodes"

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

// SetMode updates the storage mode of the full node.
func (f *FullNode) SetMode(m FullNodeMode) {
	f.Mode = m
}

// CurrentMode returns the node's current storage mode.
func (f *FullNode) CurrentMode() FullNodeMode {
	return f.Mode
}

// IsArchive reports whether the node runs in archive mode.
func (f *FullNode) IsArchive() bool {
	return f.Mode == FullNodeModeArchive
}
