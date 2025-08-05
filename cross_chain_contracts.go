package synnergy

import "sync"

// XContractMapping links a local contract address to a contract on another chain.
type XContractMapping struct {
	LocalAddress  string
	RemoteChain   string
	RemoteAddress string
}

// XContractRegistry stores cross-chain contract mappings.
type XContractRegistry struct {
	mu       sync.RWMutex
	mappings map[string]XContractMapping
}

// NewXContractRegistry creates an empty XContractRegistry.
func NewXContractRegistry() *XContractRegistry {
	return &XContractRegistry{mappings: make(map[string]XContractMapping)}
}

// RegisterMapping records a new cross-chain contract mapping.
func (r *XContractRegistry) RegisterMapping(localAddr, remoteChain, remoteAddr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.mappings[localAddr] = XContractMapping{
		LocalAddress:  localAddr,
		RemoteChain:   remoteChain,
		RemoteAddress: remoteAddr,
	}
}

// ListMappings returns all registered contract mappings.
func (r *XContractRegistry) ListMappings() []XContractMapping {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]XContractMapping, 0, len(r.mappings))
	for _, m := range r.mappings {
		out = append(out, m)
	}
	return out
}

// GetMapping retrieves a mapping by the local contract address.
func (r *XContractRegistry) GetMapping(localAddr string) (XContractMapping, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.mappings[localAddr]
	return m, ok
}

// RemoveMapping deletes a mapping by its local address.
func (r *XContractRegistry) RemoveMapping(localAddr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.mappings, localAddr)
}
