package core

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
)

// CrossChainProtocol defines a generic cross-chain integration profile.
type CrossChainProtocol struct {
	ID   string
	Name string
}

// ProtocolRegistry manages protocol definitions.
type ProtocolRegistry struct {
	mu        sync.RWMutex
	protocols map[string]*CrossChainProtocol
	nextID    uint64
}

// NewProtocolRegistry creates a ProtocolRegistry.
func NewProtocolRegistry() *ProtocolRegistry {
	return &ProtocolRegistry{protocols: make(map[string]*CrossChainProtocol)}
}

// RegisterProtocol registers a new protocol definition.
func (r *ProtocolRegistry) RegisterProtocol(name string) (*CrossChainProtocol, error) {
	if name == "" {
		return nil, errors.New("name required")
	}
	id := atomic.AddUint64(&r.nextID, 1)
	p := &CrossChainProtocol{ID: formatID("PROT", id), Name: name}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.protocols[p.ID] = p
	return p, nil
}

// ListProtocols returns all registered protocols.
func (r *ProtocolRegistry) ListProtocols() []*CrossChainProtocol {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*CrossChainProtocol, 0, len(r.protocols))
	for _, p := range r.protocols {
		res = append(res, p)
	}
	return res
}

// GetProtocol retrieves a protocol by ID.
func (r *ProtocolRegistry) GetProtocol(id string) (*CrossChainProtocol, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.protocols[id]
	return p, ok
}

// formatID creates a prefixed incremental identifier.
func formatID(prefix string, id uint64) string {
	return prefix + "-" + strconv.FormatUint(id, 10)
}
