package synnergy

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// Bridge represents a configured link between two chains.
type Bridge struct {
	ID          string
	SourceChain string
	TargetChain string
	Relayers    map[string]struct{}
}

// CrossChainManager manages bridge configurations and authorized relayers.
type CrossChainManager struct {
	mu       sync.RWMutex
	bridges  map[string]*Bridge
	relayers map[string]struct{}
}

// NewCrossChainManager creates an empty CrossChainManager instance.
func NewCrossChainManager() *CrossChainManager {
	return &CrossChainManager{
		bridges:  make(map[string]*Bridge),
		relayers: make(map[string]struct{}),
	}
}

// RegisterBridge registers a new bridge configuration and returns its ID.
func (m *CrossChainManager) RegisterBridge(sourceChain, targetChain, relayerAddr string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%d", sourceChain, targetChain, time.Now().UnixNano()))))
	bridge := &Bridge{
		ID:          id,
		SourceChain: sourceChain,
		TargetChain: targetChain,
		Relayers:    make(map[string]struct{}),
	}
	if relayerAddr != "" {
		bridge.Relayers[relayerAddr] = struct{}{}
	}
	m.bridges[id] = bridge
	return id
}

// ListBridges returns all registered bridges.
func (m *CrossChainManager) ListBridges() []*Bridge {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*Bridge, 0, len(m.bridges))
	for _, b := range m.bridges {
		out = append(out, b)
	}
	return out
}

// GetBridge retrieves a bridge by its identifier.
func (m *CrossChainManager) GetBridge(id string) (*Bridge, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.bridges[id]
	return b, ok
}

// AuthorizeRelayer whitelists a relayer address for all bridges.
func (m *CrossChainManager) AuthorizeRelayer(addr string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.relayers[addr] = struct{}{}
}

// RevokeRelayer removes a relayer from the whitelist.
func (m *CrossChainManager) RevokeRelayer(addr string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.relayers, addr)
}

// IsRelayerAuthorized returns true if the relayer is currently whitelisted.
func (m *CrossChainManager) IsRelayerAuthorized(addr string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.relayers[addr]
	return ok
}
