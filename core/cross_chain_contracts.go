package core

import (
	"errors"
	"sync"
)

// ContractMapping links a local contract address to a remote chain address.
type ContractMapping struct {
	LocalAddr   string
	RemoteChain string
	RemoteAddr  string
}

// ContractRegistry manages cross-chain contract mappings.
type ContractRegistry struct {
	mu       sync.RWMutex
	mappings map[string]*ContractMapping
}

// NewContractRegistry creates a ContractRegistry.
func NewContractRegistry() *ContractRegistry {
	return &ContractRegistry{mappings: make(map[string]*ContractMapping)}
}

// RegisterContract registers a new contract mapping.
func (r *ContractRegistry) RegisterContract(local, remoteChain, remoteAddr string) (*ContractMapping, error) {
	if local == "" || remoteChain == "" || remoteAddr == "" {
		return nil, errors.New("all fields required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.mappings[local]; exists {
		return nil, errors.New("mapping already exists")
	}
	m := &ContractMapping{LocalAddr: local, RemoteChain: remoteChain, RemoteAddr: remoteAddr}
	r.mappings[local] = m
	return m, nil
}

// GetMapping retrieves a contract mapping by local address.
func (r *ContractRegistry) GetMapping(local string) (*ContractMapping, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.mappings[local]
	return m, ok
}

// ListMappings lists all contract mappings.
func (r *ContractRegistry) ListMappings() []*ContractMapping {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*ContractMapping, 0, len(r.mappings))
	for _, m := range r.mappings {
		res = append(res, m)
	}
	return res
}

// RemoveMapping deletes a mapping by local address.
func (r *ContractRegistry) RemoveMapping(local string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.mappings[local]; !exists {
		return errors.New("mapping not found")
	}
	delete(r.mappings, local)
	return nil
}
