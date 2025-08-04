package core

import (
	"errors"
	"sync"
)

// Connection describes a link between this chain and a remote chain.
type Connection struct {
	ID          int
	LocalChain  string
	RemoteChain string
	Open        bool
}

// ConnectionManager handles connection lifecycle operations.
type ConnectionManager struct {
	mu          sync.RWMutex
	connections map[int]*Connection
	nextID      int
}

// NewConnectionManager creates a new manager.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{connections: make(map[int]*Connection)}
}

// OpenConnection establishes a new connection and returns its ID.
func (m *ConnectionManager) OpenConnection(local, remote string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	id := m.nextID
	m.connections[id] = &Connection{ID: id, LocalChain: local, RemoteChain: remote, Open: true}
	return id
}

// CloseConnection terminates an existing connection.
func (m *ConnectionManager) CloseConnection(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.connections[id]
	if !ok {
		return errors.New("connection not found")
	}
	c.Open = false
	return nil
}

// GetConnection retrieves a connection by ID.
func (m *ConnectionManager) GetConnection(id int) (*Connection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.connections[id]
	if !ok {
		return nil, errors.New("connection not found")
	}
	return c, nil
}

// ListConnections returns all known connections.
func (m *ConnectionManager) ListConnections() []*Connection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*Connection, 0, len(m.connections))
	for _, c := range m.connections {

		out = append(out, c)
	}
	return out
}
