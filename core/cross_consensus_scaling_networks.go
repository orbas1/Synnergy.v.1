package core

import (
	"errors"
	"sort"
	"sync"
)

// ConsensusNetwork represents a connection between differing consensus systems.
type ConsensusNetwork struct {
	ID              int
	SourceConsensus string
	TargetConsensus string
}

// ConsensusNetworkManager manages registered consensus networks.
type ConsensusNetworkManager struct {
	mu       sync.RWMutex
	networks map[int]ConsensusNetwork
	nextID   int
	relayers map[string]bool
}

// NewConsensusNetworkManager creates a new manager.
func NewConsensusNetworkManager() *ConsensusNetworkManager {
	return &ConsensusNetworkManager{networks: make(map[int]ConsensusNetwork), relayers: make(map[string]bool)}
}

// AuthorizeRelayer whitelists an address for network modifications.
func (m *ConsensusNetworkManager) AuthorizeRelayer(addr string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.relayers[addr] = true
}

// RevokeRelayer removes an address from the whitelist.
func (m *ConsensusNetworkManager) RevokeRelayer(addr string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.relayers, addr)
}

// IsRelayerAuthorized returns true if the address is authorized.
func (m *ConsensusNetworkManager) IsRelayerAuthorized(addr string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.relayers[addr]
}

// RegisterNetwork registers a new cross-consensus network and returns its ID.
// Only authorized relayers may register networks.
func (m *ConsensusNetworkManager) RegisterNetwork(source, target, relayer string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.relayers[relayer] {
		return 0, errors.New("unauthorized relayer")
	}
	m.nextID++
	id := m.nextID
	m.networks[id] = ConsensusNetwork{ID: id, SourceConsensus: source, TargetConsensus: target}
	return id, nil
}

// RemoveNetwork deletes a registered network by ID if the relayer is authorized.
func (m *ConsensusNetworkManager) RemoveNetwork(id int, relayer string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.relayers[relayer] {
		return errors.New("unauthorized relayer")
	}
	if _, ok := m.networks[id]; !ok {
		return errors.New("network not found")
	}
	delete(m.networks, id)
	return nil
}

// ListNetworks returns all registered networks.
func (m *ConsensusNetworkManager) ListNetworks() []ConsensusNetwork {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]ConsensusNetwork, 0, len(m.networks))
	for _, n := range m.networks {
		out = append(out, n)
	}
	return out
}

// AuthorizedRelayers returns a sorted list of whitelisted relayer addresses.
// The slice is safe for reuse by callers.
func (m *ConsensusNetworkManager) AuthorizedRelayers() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	relayers := make([]string, 0, len(m.relayers))
	for addr := range m.relayers {
		relayers = append(relayers, addr)
	}
	sort.Strings(relayers)
	return relayers
}

// GetNetwork retrieves a network configuration by ID.
func (m *ConsensusNetworkManager) GetNetwork(id int) (ConsensusNetwork, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n, ok := m.networks[id]
	if !ok {
		return ConsensusNetwork{}, errors.New("network not found")
	}
	return n, nil
}
