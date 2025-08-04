package core

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
)

// CCSNetwork represents a bridge between two independent consensus systems.
type CCSNetwork struct {
	ID              string
	SourceConsensus string
	TargetConsensus string
}

// CCSNetworkRegistry manages cross-consensus networks.
type CCSNetworkRegistry struct {
	mu       sync.RWMutex
	networks map[string]*CCSNetwork
	nextID   uint64
}

// NewCCSNetworkRegistry creates a CCSNetworkRegistry.
func NewCCSNetworkRegistry() *CCSNetworkRegistry {
	return &CCSNetworkRegistry{networks: make(map[string]*CCSNetwork)}
}

// RegisterNetwork registers a cross-consensus network.
func (r *CCSNetworkRegistry) RegisterNetwork(source, target string) (*CCSNetwork, error) {
	if source == "" || target == "" {
		return nil, errors.New("consensus names required")
	}
	id := atomic.AddUint64(&r.nextID, 1)
	n := &CCSNetwork{ID: formatID("CCS", id), SourceConsensus: source, TargetConsensus: target}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.networks[n.ID] = n
	return n, nil
}

// ListNetworks lists configured networks.
func (r *CCSNetworkRegistry) ListNetworks() []*CCSNetwork {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*CCSNetwork, 0, len(r.networks))
	for _, n := range r.networks {
		res = append(res, n)
	}
	return res
}

// GetNetwork retrieves a network configuration.
func (r *CCSNetworkRegistry) GetNetwork(id string) (*CCSNetwork, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.networks[id]
	return n, ok
}

// formatID creates a prefixed incremental identifier.
func formatID(prefix string, id uint64) string {
	return prefix + "-" + strconv.FormatUint(id, 10)
}
