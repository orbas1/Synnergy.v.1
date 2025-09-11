package core

import (
	"errors"
	"sync"

	"synnergy/internal/nodes"
)

// LightNode maintains a minimal view of the chain using block headers only.
type LightNode struct {
	*BaseNode
	mu      sync.RWMutex
	headers []nodes.BlockHeader
}

// NewLightNode constructs a light node with no headers.
func NewLightNode(id nodes.Address) *LightNode {
	return &LightNode{BaseNode: NewBaseNode(id)}
}

// AddHeader appends a header to the local view.
func (n *LightNode) AddHeader(h nodes.BlockHeader) error {
	if h.Hash == "" {
		return errors.New("empty header hash")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.headers = append(n.headers, h)
	return nil
}

// LatestHeader returns the most recently added header.
func (n *LightNode) LatestHeader() (nodes.BlockHeader, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if len(n.headers) == 0 {
		return nodes.BlockHeader{}, false
	}
	return n.headers[len(n.headers)-1], true
}

// Headers returns a copy of all stored block headers.
func (n *LightNode) Headers() []nodes.BlockHeader {
	n.mu.RLock()
	defer n.mu.RUnlock()
	out := make([]nodes.BlockHeader, len(n.headers))
	copy(out, n.headers)
	return out
}
