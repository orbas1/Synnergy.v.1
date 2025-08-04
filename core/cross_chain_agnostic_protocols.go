package core

import "sync"

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
}

// NewProtocolRegistry creates a new registry instance.
func NewProtocolRegistry() *ProtocolRegistry {
	return &ProtocolRegistry{protocols: make(map[int]ProtocolDefinition)}
}

// Register adds a new protocol by name and returns its identifier.
func (r *ProtocolRegistry) Register(name string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	id := r.nextID
	r.protocols[id] = ProtocolDefinition{ID: id, Name: name}
	return id
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
