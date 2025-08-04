package core

import (
	"fmt"
	"sync"
)

// ChainConnection represents an active cross-chain connection between two chains.
type ChainConnection struct {
	ID          string
	LocalChain  string
	RemoteChain string
	Active      bool
}

// ConnectionRegistry tracks cross-chain connections.
type ConnectionRegistry struct {
	mu          sync.RWMutex
	seq         int
	connections map[string]*ChainConnection
}

// NewConnectionRegistry creates a new registry.
func NewConnectionRegistry() *ConnectionRegistry {
	return &ConnectionRegistry{connections: make(map[string]*ChainConnection)}
}

// OpenConnection establishes a new connection.
func (r *ConnectionRegistry) OpenConnection(local, remote string) (*ChainConnection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	id := fmt.Sprintf("conn-%d", r.seq)
	c := &ChainConnection{ID: id, LocalChain: local, RemoteChain: remote, Active: true}
	r.connections[id] = c
	return c, nil
}

// CloseConnection terminates a connection.
func (r *ConnectionRegistry) CloseConnection(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.connections[id]
	if !ok {
		return fmt.Errorf("connection %s not found", id)
	}
	c.Active = false
	return nil
}

// GetConnection retrieves connection details.
func (r *ConnectionRegistry) GetConnection(id string) (*ChainConnection, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.connections[id]
	return c, ok
}

// ListConnections lists all connections.
func (r *ConnectionRegistry) ListConnections() []*ChainConnection {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*ChainConnection, 0, len(r.connections))
	for _, c := range r.connections {
		out = append(out, c)
	}
	return out
}
