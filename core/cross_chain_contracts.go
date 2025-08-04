package core

import "sync"

// ContractMapping links a local contract address to a remote chain address.
type ContractMapping struct {
	LocalAddr   string
	RemoteChain string
	RemoteAddr  string
}

// CrossChainRegistry manages cross-chain contract mappings.
type CrossChainRegistry struct {
	mu       sync.RWMutex
	mappings map[string]*ContractMapping
}

// NewCrossChainRegistry creates an empty registry.
func NewCrossChainRegistry() *CrossChainRegistry {
	return &CrossChainRegistry{mappings: make(map[string]*ContractMapping)}
}

// RegisterMapping registers a new contract mapping.
func (r *CrossChainRegistry) RegisterMapping(local, remoteChain, remoteAddr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.mappings[local] = &ContractMapping{LocalAddr: local, RemoteChain: remoteChain, RemoteAddr: remoteAddr}
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
func (r *CrossChainRegistry) RemoveMapping(local string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.mappings, local)
}
