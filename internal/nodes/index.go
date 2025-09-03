package nodes

import (
	"errors"
	"sync"
)

// Node defines minimal behaviour required by all network nodes.  The
// interface is intentionally small so that specialised implementations can
// embed additional functionality without pulling in dependencies from the
// `core` packages.
type Node interface {
	// ID returns the node identifier.
	ID() Address
	// Start begins node operations such as networking routines.
	Start() error
	// Stop gracefully halts node operations.
	Stop() error
	// IsRunning reports whether the node is currently active.
	IsRunning() bool
	// Peers returns identifiers for all known peers.
	Peers() []Address
	// DialSeed connects the node to a seed peer by address.
	DialSeed(addr Address) error
}

// NodeInterface is kept for backward compatibility with earlier stages.
// New code should depend on Node directly.
type NodeInterface = Node

// BasicNode provides a concurrency safe reference implementation of the
// Node interface.  Concrete node types can embed this struct to obtain a
// consistent lifecycle and peer management behaviour.
type BasicNode struct {
	id      Address
	mu      sync.RWMutex
	running bool
	peers   []Address
}

// NewBasicNode constructs a new BasicNode with the given identifier.
func NewBasicNode(id Address) *BasicNode {
	return &BasicNode{id: id}
}

// ID returns the identifier of the node.
func (n *BasicNode) ID() Address {
	return n.id
}

// Start marks the node as running.  It returns an error if the node is
// already active.
func (n *BasicNode) Start() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.running {
		return errors.New("node already running")
	}
	n.running = true
	return nil
}

// Stop marks the node as stopped.  It returns an error if the node was not
// previously running.
func (n *BasicNode) Stop() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.running {
		return errors.New("node not running")
	}
	n.running = false
	return nil
}

// IsRunning reports whether the node is currently active.
func (n *BasicNode) IsRunning() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.running
}

// Peers returns a copy of the peer list known to the node.
func (n *BasicNode) Peers() []Address {
	n.mu.RLock()
	defer n.mu.RUnlock()
	cp := make([]Address, len(n.peers))
	copy(cp, n.peers)
	return cp
}

// DialSeed records a connection to a seed peer.  The implementation avoids
// duplicates but otherwise does not attempt to establish a real network
// connection.
func (n *BasicNode) DialSeed(addr Address) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, p := range n.peers {
		if p == addr {
			return nil
		}
	}
	n.peers = append(n.peers, addr)
	return nil
}
