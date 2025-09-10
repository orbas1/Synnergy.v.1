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

// CrossChainRegistry manages cross-chain contract mappings.
type CrossChainRegistry struct {
	mu       sync.RWMutex
	mappings map[string]*ContractMapping
	relayers map[string]bool
}

// NewCrossChainRegistry creates an empty registry.
func NewCrossChainRegistry() *CrossChainRegistry {
	return &CrossChainRegistry{
		mappings: make(map[string]*ContractMapping),
		relayers: make(map[string]bool),
	}
}

// AuthorizeRelayer adds a relayer to the whitelist.
func (r *CrossChainRegistry) AuthorizeRelayer(relayer string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.relayers[relayer] = true
}

// RevokeRelayer removes a relayer from the whitelist.
func (r *CrossChainRegistry) RevokeRelayer(relayer string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.relayers, relayer)
}

// IsRelayerAuthorized checks if a relayer is authorized.
func (r *CrossChainRegistry) IsRelayerAuthorized(relayer string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.relayers[relayer]
}

// RegisterMapping registers a new contract mapping. Only authorized relayers may register.
func (r *CrossChainRegistry) RegisterMapping(relayer, local, remoteChain, remoteAddr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.relayers[relayer] {
		return errors.New("unauthorized relayer")
	}
	r.mappings[local] = &ContractMapping{LocalAddress: local, RemoteChain: remoteChain, RemoteAddress: remoteAddr}
	return nil
}

// GetMapping retrieves a mapping by local address.
func (r *CrossChainRegistry) GetMapping(local string) (*ContractMapping, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.mappings[local]
	return m, ok
}

// ListMappings returns all registered mappings.
func (r *CrossChainRegistry) ListMappings() []*ContractMapping {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*ContractMapping, 0, len(r.mappings))
	for _, m := range r.mappings {
		out = append(out, m)
	}
	return out
}

// RemoveMapping deletes a mapping by local address.
func (r *CrossChainRegistry) RemoveMapping(relayer, local string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.relayers[relayer] {
		return errors.New("unauthorized relayer")
	}
	if _, ok := r.mappings[local]; !ok {
		return errors.New("mapping not found")
	}
	delete(r.mappings, local)
	return nil
}
