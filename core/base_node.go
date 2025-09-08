package core

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"sync"

	"synnergy/internal/nodes"
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

// IsRunning reports whether the node is currently active.
func (n *BaseNode) IsRunning() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.running
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

// DialSeedSigned records a connection to a seed peer after verifying the
// provided signature matches the peer's address.
func (n *BaseNode) DialSeedSigned(addr nodes.Address, sig []byte, pub ed25519.PublicKey) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.running {
		return fmt.Errorf("node not running")
	}
	if hex.EncodeToString(pub) != string(addr) {
		return fmt.Errorf("address mismatch")
	}
	if !ed25519.Verify(pub, []byte(addr), sig) {
		return fmt.Errorf("invalid signature")
	}
	n.peers[addr] = struct{}{}
	return nil
}

var _ nodes.NodeInterface = (*BaseNode)(nil)
