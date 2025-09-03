package nodes

import (
	"sync"

	"synnergy"
)

// HolographicNode provides holographic data distribution and redundancy.
type HolographicNode struct {
	*BasicNode
	mu    sync.RWMutex
	store map[string]synnergy.HolographicFrame
}

// NewHolographicNode creates a new HolographicNode with the given identifier.
func NewHolographicNode(id Address) *HolographicNode {
	return &HolographicNode{
		BasicNode: NewBasicNode(id),
		store:     make(map[string]synnergy.HolographicFrame),
	}
}

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
