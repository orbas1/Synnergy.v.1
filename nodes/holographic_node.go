package nodes

import (
	"sync"

	"synnergy"
)

// HolographicNode provides holographic data distribution and redundancy.
type HolographicNode struct {
	id    Address
	mu    sync.RWMutex
	store map[string]synnergy.HolographicFrame
	peers map[Address]struct{}
}

// NewHolographicNode creates a new HolographicNode with the given identifier.
func NewHolographicNode(id Address) *HolographicNode {
	return &HolographicNode{
		id:    id,
		store: make(map[string]synnergy.HolographicFrame),
		peers: make(map[Address]struct{}),
	}
}

// ID returns the node identifier.
func (n *HolographicNode) ID() Address { return n.id }

// Start implements the NodeInterface; holographic nodes currently have no
// background processes so Start is a no-op.
func (n *HolographicNode) Start() error { return nil }

// Stop implements the NodeInterface; holographic nodes currently have no
// background processes so Stop is a no-op.
func (n *HolographicNode) Stop() error { return nil }

// Peers returns all known peer addresses.
func (n *HolographicNode) Peers() []Address {
	n.mu.RLock()
	defer n.mu.RUnlock()
	out := make([]Address, 0, len(n.peers))
	for p := range n.peers {
		out = append(out, p)
	}
	return out
}

// DialSeed adds the provided address to the peer list.
func (n *HolographicNode) DialSeed(addr Address) error {
	n.mu.Lock()
	n.peers[addr] = struct{}{}
	n.mu.Unlock()
	return nil
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
