package core

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
)

// Bridge defines parameters for a cross-chain bridge.
type Bridge struct {
	ID          string
	SourceChain string
	TargetChain string
	Relayers    map[string]bool
}

// BridgeRegistry manages cross-chain bridges.
type BridgeRegistry struct {
	mu      sync.RWMutex
	bridges map[string]*Bridge
	nextID  uint64
}

// NewBridgeRegistry creates an empty BridgeRegistry.
func NewBridgeRegistry() *BridgeRegistry {
	return &BridgeRegistry{bridges: make(map[string]*Bridge)}
}

// RegisterBridge registers a new bridge between two chains and whitelists an initial relayer.
func (r *BridgeRegistry) RegisterBridge(source, target, relayer string) (*Bridge, error) {
	if source == "" || target == "" {
		return nil, errors.New("source and target chains required")
	}
	id := atomic.AddUint64(&r.nextID, 1)
	bridge := &Bridge{
		ID:          formatID("BRG", id),
		SourceChain: source,
		TargetChain: target,
		Relayers:    map[string]bool{relayer: true},
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bridges[bridge.ID] = bridge
	return bridge, nil
}

// ListBridges returns all registered bridges.
func (r *BridgeRegistry) ListBridges() []*Bridge {
	r.mu.RLock()
	defer r.mu.RUnlock()
	bridges := make([]*Bridge, 0, len(r.bridges))
	for _, b := range r.bridges {
		bridges = append(bridges, b)
	}
	return bridges
}

// GetBridge retrieves a bridge by ID.
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
		return errors.New("bridge not found")
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
		return errors.New("bridge not found")
	}
	delete(b.Relayers, relayer)
	return nil
}

// formatID creates a prefixed incremental identifier.
func formatID(prefix string, id uint64) string {
	return prefix + "-" + strconv.FormatUint(id, 10)
}
