package nodes

import (
	"sync"

	"synnergy"
)

// HolographicNode provides holographic data distribution and redundancy.
type HolographicNode struct {
	id    string
	mu    sync.RWMutex
	store map[string]synnergy.HolographicFrame
}

// NewHolographicNode creates a new HolographicNode with the given identifier.
func NewHolographicNode(id string) *HolographicNode {
	return &HolographicNode{
		id:    id,
		store: make(map[string]synnergy.HolographicFrame),
	}
}

// ID returns the node identifier.
func (n *HolographicNode) ID() string { return n.id }

// Start implements the NodeInterface; holographic nodes currently have no
// background processes so Start is a no-op.
func (n *HolographicNode) Start() error { return nil }

// Stop implements the NodeInterface; holographic nodes currently have no
// background processes so Stop is a no-op.
func (n *HolographicNode) Stop() error { return nil }

// Store saves a holographic frame in the node's internal storage.
func (n *HolographicNode) Store(frame synnergy.HolographicFrame) {
	n.mu.Lock()
	n.store[frame.ID] = frame
	n.mu.Unlock()
}

// Retrieve fetches a holographic frame by ID. The returned boolean indicates
// whether the frame was found.
func (n *HolographicNode) Retrieve(id string) (synnergy.HolographicFrame, bool) {
	n.mu.RLock()
	frame, ok := n.store[id]
	n.mu.RUnlock()
	return frame, ok
}
