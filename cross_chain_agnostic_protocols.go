package synnergy

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// CrossChainProtocol defines a protocol standard understood across chains.
type CrossChainProtocol struct {
	ID   string
	Name string
}

// ProtocolRegistry stores registered protocol definitions.
type ProtocolRegistry struct {
	mu        sync.RWMutex
	protocols map[string]CrossChainProtocol
}

// NewProtocolRegistry creates an empty ProtocolRegistry.
func NewProtocolRegistry() *ProtocolRegistry {
	return &ProtocolRegistry{protocols: make(map[string]CrossChainProtocol)}
}

// RegisterProtocol registers a new protocol and returns its identifier.
func (r *ProtocolRegistry) RegisterProtocol(name string) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s|%d", name, time.Now().UnixNano()))))
	r.protocols[id] = CrossChainProtocol{ID: id, Name: name}
	return id
}

// ListProtocols returns all registered protocols.
func (r *ProtocolRegistry) ListProtocols() []CrossChainProtocol {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]CrossChainProtocol, 0, len(r.protocols))
	for _, p := range r.protocols {
		out = append(out, p)
	}
	return out
}

// GetProtocol retrieves a protocol definition by ID.
func (r *ProtocolRegistry) GetProtocol(id string) (CrossChainProtocol, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.protocols[id]
	return p, ok
}
