package core

import (
	"errors"
	"sync"
)

// ConsensusNetwork represents a connection between differing consensus systems.
type ConsensusNetwork struct {
	ID              int


// ConsensusNetworkManager manages registered consensus networks.
type ConsensusNetworkManager struct {
	mu       sync.RWMutex
	networks map[int]ConsensusNetwork
	nextID   int
}

// NewConsensusNetworkManager creates a new manager.
func NewConsensusNetworkManager() *ConsensusNetworkManager {
	return &ConsensusNetworkManager{networks: make(map[int]ConsensusNetwork)}
}

// RegisterNetwork registers a new cross-consensus network and returns its ID.
func (m *ConsensusNetworkManager) RegisterNetwork(source, target string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	id := m.nextID
	m.networks[id] = ConsensusNetwork{ID: id, SourceConsensus: source, TargetConsensus: target}
	return id
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
