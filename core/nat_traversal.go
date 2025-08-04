package core

import "sync"

// NATManager tracks port mappings for nodes operating behind NAT devices. It is
// a lightweight helper to register and discover externally reachable ports.
type NATManager struct {
	mu       sync.RWMutex
	mappings map[string]int // node id -> external port
}

// NewNATManager creates an empty NAT manager.
func NewNATManager() *NATManager {
	return &NATManager{mappings: make(map[string]int)}
}

// MapPort records an external port for a node.
func (n *NATManager) MapPort(id string, port int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.mappings[id] = port
}

// GetPort retrieves the mapped port for a node if available.
func (n *NATManager) GetPort(id string) (int, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	p, ok := n.mappings[id]
	return p, ok
}

// RemoveMapping deletes a node's port mapping.
func (n *NATManager) RemoveMapping(id string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.mappings, id)
}
