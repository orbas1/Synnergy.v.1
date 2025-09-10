package core

import (
	"fmt"
	"sync"
)

// Bridge defines parameters for a cross-chain bridge.
type Bridge struct {
	ID          string
	SourceChain string
	TargetChain string
	Relayers    map[string]bool
}

// BridgeRegistry manages bridges and authorized relayers.
type BridgeRegistry struct {
	mu      sync.RWMutex
	seq     int
	bridges map[string]*Bridge
}

// NewBridgeRegistry creates an empty BridgeRegistry.
func NewBridgeRegistry() *BridgeRegistry {
	return &BridgeRegistry{bridges: make(map[string]*Bridge)}
}

// RegisterBridge registers a new bridge and authorizes the initial relayer.
func (r *BridgeRegistry) RegisterBridge(source, target, relayer string) (*Bridge, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	id := fmt.Sprintf("bridge-%d", r.seq)
	b := &Bridge{
		ID:          id,
		SourceChain: source,
		TargetChain: target,
		Relayers:    map[string]bool{relayer: true},
	}
	r.bridges[id] = b
	return b, nil
}

// ListBridges returns all registered bridges.
func (r *BridgeRegistry) ListBridges() []*Bridge {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Bridge, 0, len(r.bridges))
	for _, b := range r.bridges {
		out = append(out, b)
	}
	return out
}

// GetBridge retrieves a bridge by its ID.
func (r *BridgeRegistry) GetBridge(id string) (*Bridge, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.bridges[id]
	return b, ok
}

// AuthorizeRelayer whitelists a relayer for a bridge.
func (r *BridgeRegistry) AuthorizeRelayer(id, relayer string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.bridges[id]
	if !ok {
		return fmt.Errorf("bridge %s not found", id)
	}
	if b.Relayers == nil {
		b.Relayers = make(map[string]bool)
	}
	b.Relayers[relayer] = true
	return nil
}

// RevokeRelayer removes a relayer from the whitelist.
func (r *BridgeRegistry) RevokeRelayer(id, relayer string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.bridges[id]
	if !ok {
		return fmt.Errorf("bridge %s not found", id)
	}
	delete(b.Relayers, relayer)
	return nil
}

// IsRelayerAuthorized checks if a relayer is authorized for a bridge.
// It returns true if the bridge exists and the relayer is whitelisted.
func (r *BridgeRegistry) IsRelayerAuthorized(id, relayer string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.bridges[id]
	if !ok {
		return false
	}
	return b.Relayers[relayer]
}

// RemoveBridge deletes a bridge from the registry.
// It returns an error if the bridge cannot be found.
func (r *BridgeRegistry) RemoveBridge(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.bridges[id]; !ok {
		return fmt.Errorf("bridge %s not found", id)
	}
	delete(r.bridges, id)
	return nil
}
