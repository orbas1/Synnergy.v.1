package core

import (
	"fmt"
	"sync"

	"synnergy/nodes"
)

// BaseNode wraps a NodeInterface and exposes common networking behaviour.
type BaseNode struct {
	id      nodes.Address
	peers   map[nodes.Address]struct{}
	running bool
	mu      sync.RWMutex
}

// NewBaseNode constructs a BaseNode with the provided identifier.
func NewBaseNode(id nodes.Address) *BaseNode {
	return &BaseNode{
		id:    id,
		peers: make(map[nodes.Address]struct{}),
	}
}

// ID returns the node identifier.
func (n *BaseNode) ID() nodes.Address { return n.id }

// Start marks the node as running.
func (n *BaseNode) Start() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.running {
		return nil
	}
	n.running = true
	return nil
}

// Stop halts node operations.
func (n *BaseNode) Stop() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.running {
		return nil
	}
	n.running = false
	return nil
}

// Peers returns the list of known peers.
func (n *BaseNode) Peers() []nodes.Address {
	n.mu.RLock()
	defer n.mu.RUnlock()
	out := make([]nodes.Address, 0, len(n.peers))
	for p := range n.peers {
		out = append(out, p)
	}
	return out
}

// DialSeed records a connection to a seed peer.
func (n *BaseNode) DialSeed(addr nodes.Address) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.running {
		return fmt.Errorf("node not running")
	}
	n.peers[addr] = struct{}{}
	return nil
}

var _ nodes.NodeInterface = (*BaseNode)(nil)
