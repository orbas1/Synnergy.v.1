package core

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
)

// ChainConnection represents an active cross-chain connection between two chains.
type ChainConnection struct {
	ID          string
	LocalChain  string
	RemoteChain string
	Open        bool
}

// ConnectionManager manages cross-chain connections.
type ConnectionManager struct {
	mu          sync.RWMutex
	connections map[string]*ChainConnection
	nextID      uint64
}

// NewConnectionManager creates a ConnectionManager.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{connections: make(map[string]*ChainConnection)}
}

// OpenConnection establishes a new connection.
func (m *ConnectionManager) OpenConnection(local, remote string) (*ChainConnection, error) {
	if local == "" || remote == "" {
		return nil, errors.New("chain names required")
	}
	id := atomic.AddUint64(&m.nextID, 1)
	c := &ChainConnection{ID: formatID("CON", id), LocalChain: local, RemoteChain: remote, Open: true}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[c.ID] = c
	return c, nil
}

// CloseConnection terminates a connection.
func (m *ConnectionManager) CloseConnection(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.connections[id]
	if !ok {
		return errors.New("connection not found")
	}
	c.Open = false
	return nil
}

// GetConnection retrieves a connection.
func (m *ConnectionManager) GetConnection(id string) (*ChainConnection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.connections[id]
	return c, ok
}

// ListConnections lists all connections.
func (m *ConnectionManager) ListConnections() []*ChainConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]*ChainConnection, 0, len(m.connections))
	for _, c := range m.connections {
		res = append(res, c)
	}
	return res
}

// formatID creates a prefixed incremental identifier.
func formatID(prefix string, id uint64) string {
	return prefix + "-" + strconv.FormatUint(id, 10)
}
