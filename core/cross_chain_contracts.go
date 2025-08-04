package core

import (
	"errors"
	"sync"
)

// ContractMapping links a local contract to a remote chain address.
type ContractMapping struct {
	LocalAddress  string
	RemoteChain   string
	RemoteAddress string
}

// ContractRegistry manages contract mappings for cross-chain calls.
type ContractRegistry struct {
	mu       sync.RWMutex
	mappings map[string]ContractMapping
}

// NewContractRegistry creates a new registry.
func NewContractRegistry() *ContractRegistry {
	return &ContractRegistry{mappings: make(map[string]ContractMapping)}
}

// RegisterMapping stores a new mapping in the registry.
func (r *ContractRegistry) RegisterMapping(localAddr, remoteChain, remoteAddr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.mappings[localAddr] = ContractMapping{LocalAddress: localAddr, RemoteChain: remoteChain, RemoteAddress: remoteAddr}
}

// GetMapping retrieves mapping information for a local contract address.
func (r *ContractRegistry) GetMapping(localAddr string) (ContractMapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.mappings[localAddr]
	if !ok {
		return ContractMapping{}, errors.New("mapping not found")
	}
	return m, nil
}

// ListMappings returns all registered mappings.
func (r *ContractRegistry) ListMappings() []ContractMapping {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]ContractMapping, 0, len(r.mappings))
	for _, m := range r.mappings {
		out = append(out, m)
	}
	return out
}

// RemoveMapping deletes a mapping from the registry.
func (r *ContractRegistry) RemoveMapping(localAddr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.mappings[localAddr]; !ok {
		return errors.New("mapping not found")
	}
	delete(r.mappings, localAddr)
	return nil

}
