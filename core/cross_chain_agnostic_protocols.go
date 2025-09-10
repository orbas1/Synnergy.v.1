package core

import (
	"errors"
	"sync"
)

// ProtocolDefinition represents a cross-chain protocol definition.
type ProtocolDefinition struct {
	ID   int
	Name string
}

// ProtocolRegistry manages protocol definitions used for cross-chain messaging.
type ProtocolRegistry struct {
	mu        sync.RWMutex
	protocols map[int]ProtocolDefinition
	nextID    int
	relayers  map[string]bool
}

// NewProtocolRegistry creates a new registry instance.
func NewProtocolRegistry() *ProtocolRegistry {
	return &ProtocolRegistry{
		protocols: make(map[int]ProtocolDefinition),
		relayers:  make(map[string]bool),
	}
}

// AuthorizeRelayer adds a relayer to the whitelist.
func (r *ProtocolRegistry) AuthorizeRelayer(relayer string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.relayers[relayer] = true
}

// RevokeRelayer removes a relayer from the whitelist.
func (r *ProtocolRegistry) RevokeRelayer(relayer string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.relayers, relayer)
}

// IsRelayerAuthorized checks if a relayer is authorized.
func (r *ProtocolRegistry) IsRelayerAuthorized(relayer string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.relayers[relayer]
}

// Register adds a new protocol by name and returns its identifier.
// Only authorized relayers may register new protocols.
func (r *ProtocolRegistry) Register(name, relayer string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.relayers[relayer] {
		return 0, errors.New("unauthorized relayer")
	}
	r.nextID++
	id := r.nextID
	r.protocols[id] = ProtocolDefinition{ID: id, Name: name}
	return id, nil
}

// Remove deletes a protocol definition by identifier if authorized.
func (r *ProtocolRegistry) Remove(id int, relayer string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.relayers[relayer] {
		return errors.New("unauthorized relayer")
	}
	if _, ok := r.protocols[id]; !ok {
		return errors.New("protocol not found")
	}
	delete(r.protocols, id)
	return nil
}

// List returns all registered protocol definitions.
func (r *ProtocolRegistry) List() []ProtocolDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]ProtocolDefinition, 0, len(r.protocols))
	for _, p := range r.protocols {
		out = append(out, p)
	}
	return out
}

// Get retrieves a protocol definition by identifier.
func (r *ProtocolRegistry) Get(id int) (ProtocolDefinition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.protocols[id]
	return p, ok
}
