package core

import "sync"

// ContractMapping links a local contract address to a remote chain address.
type ContractMapping struct {
	LocalAddr   string
	RemoteChain string
	RemoteAddr  string
}

// CrossChainContractRegistry manages cross-chain contract mappings.
type CrossChainContractRegistry struct {
        mu       sync.RWMutex
        mappings map[string]*ContractMapping
}

// NewCrossChainContractRegistry creates an empty registry.
func NewCrossChainContractRegistry() *CrossChainContractRegistry {
        return &CrossChainContractRegistry{mappings: make(map[string]*ContractMapping)}
}

// RegisterMapping registers a new contract mapping.
func (r *CrossChainContractRegistry) RegisterMapping(local, remoteChain, remoteAddr string) {
        r.mu.Lock()
        defer r.mu.Unlock()
        r.mappings[local] = &ContractMapping{LocalAddr: local, RemoteChain: remoteChain, RemoteAddr: remoteAddr}
}

// GetMapping retrieves a mapping by local address.
func (r *CrossChainContractRegistry) GetMapping(local string) (*ContractMapping, bool) {
        r.mu.RLock()
        defer r.mu.RUnlock()
        m, ok := r.mappings[local]
        return m, ok
}

// ListMappings returns all registered mappings.
func (r *CrossChainContractRegistry) ListMappings() []*ContractMapping {
        r.mu.RLock()
        defer r.mu.RUnlock()
        out := make([]*ContractMapping, 0, len(r.mappings))
        for _, m := range r.mappings {
                out = append(out, m)
        }
        return out
}

// RemoveMapping deletes a mapping by local address.
func (r *CrossChainContractRegistry) RemoveMapping(local string) {
        r.mu.Lock()
        defer r.mu.Unlock()
        delete(r.mappings, local)
}
