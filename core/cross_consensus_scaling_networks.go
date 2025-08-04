package core

import (
	"fmt"
	"sync"
)

// CCSNetwork represents a bridge between two independent consensus systems.
type CCSNetwork struct {
	ID              string
	SourceConsensus string
	TargetConsensus string
}

// CCSRegistry manages cross-consensus networks.
type CCSRegistry struct {
	mu       sync.RWMutex
	seq      int
	networks map[string]*CCSNetwork
}

// NewCCSRegistry creates a new registry.
func NewCCSRegistry() *CCSRegistry {
	return &CCSRegistry{networks: make(map[string]*CCSNetwork)}
}

// RegisterNetwork registers a cross-consensus network.
func (r *CCSRegistry) RegisterNetwork(source, target string) *CCSNetwork {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	id := fmt.Sprintf("ccsn-%d", r.seq)
	n := &CCSNetwork{ID: id, SourceConsensus: source, TargetConsensus: target}
	r.networks[id] = n
	return n
}

// GetNetwork retrieves a network configuration by ID.
func (r *CCSRegistry) GetNetwork(id string) (*CCSNetwork, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.networks[id]
	return n, ok
}

// ListNetworks lists all configured networks.
func (r *CCSRegistry) ListNetworks() []*CCSNetwork {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*CCSNetwork, 0, len(r.networks))
	for _, n := range r.networks {
		out = append(out, n)
	}
	return out
}
