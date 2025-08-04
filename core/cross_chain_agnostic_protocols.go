package core

import (
	"fmt"
	"sync"
)

// CrossChainProtocol defines a generic cross-chain integration profile.
type CrossChainProtocol struct {
	ID   string
	Name string
}

// ProtocolRegistry stores registered cross-chain protocols.
type ProtocolRegistry struct {
	mu        sync.RWMutex
	seq       int
	protocols map[string]*CrossChainProtocol
}

// NewProtocolRegistry creates an empty registry.
func NewProtocolRegistry() *ProtocolRegistry {
	return &ProtocolRegistry{protocols: make(map[string]*CrossChainProtocol)}
}

// RegisterProtocol adds a new protocol definition.
func (r *ProtocolRegistry) RegisterProtocol(name string) *CrossChainProtocol {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	id := fmt.Sprintf("proto-%d", r.seq)
	p := &CrossChainProtocol{ID: id, Name: name}
	r.protocols[id] = p
	return p
}

// GetProtocol retrieves a protocol by ID.
func (r *ProtocolRegistry) GetProtocol(id string) (*CrossChainProtocol, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.protocols[id]
	return p, ok
}

// ListProtocols lists all registered protocols.
func (r *ProtocolRegistry) ListProtocols() []*CrossChainProtocol {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*CrossChainProtocol, 0, len(r.protocols))
	for _, p := range r.protocols {
		out = append(out, p)
	}
	return out
}
