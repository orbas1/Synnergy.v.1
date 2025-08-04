package core

import (
	"errors"
	"sync"
)

// ChainConnection describes a link between this chain and a remote chain.
type ChainConnection struct {
	ID          int
	LocalChain  string
	RemoteChain string
	Open        bool
}

// ChainConnectionManager handles connection lifecycle operations.
type ChainConnectionManager struct {
	mu          sync.RWMutex
	connections map[int]*ChainConnection
	nextID      int
}

// NewChainConnectionManager creates a new manager.
func NewChainConnectionManager() *ChainConnectionManager {
	return &ChainConnectionManager{connections: make(map[int]*ChainConnection)}
}

// Open establishes a new connection and returns its ID.
func (m *ChainConnectionManager) Open(local, remote string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	id := m.nextID
	m.connections[id] = &ChainConnection{ID: id, LocalChain: local, RemoteChain: remote, Open: true}
	return id
}

// Close terminates an existing connection.
func (m *ChainConnectionManager) Close(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.connections[id]
	if !ok {
		return errors.New("connection not found")
	}
	c.Open = false
	return nil
}

// Get retrieves a connection by ID.
func (m *ChainConnectionManager) Get(id int) (*ChainConnection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.connections[id]
	if !ok {
		return nil, errors.New("connection not found")
	}
	return c, nil
}

// List returns all known connections.
func (m *ChainConnectionManager) List() []*ChainConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*ChainConnection, 0, len(m.connections))
	for _, c := range m.connections {
		out = append(out, c)
	}
	return out
}
