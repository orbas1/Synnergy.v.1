package synnergy

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// ChainConnection represents an active or historic cross-chain connection.
type ChainConnection struct {
	ID          string
	LocalChain  string
	RemoteChain string
	OpenedAt    time.Time
	ClosedAt    time.Time
	Closed      bool
}

// ConnectionManager manages cross-chain connections.
type ConnectionManager struct {
	mu          sync.RWMutex
	connections map[string]*ChainConnection
}

// NewConnectionManager creates an empty ConnectionManager.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{connections: make(map[string]*ChainConnection)}
}

// OpenConnection establishes a new link between two chains and returns its ID.
func (m *ConnectionManager) OpenConnection(localChain, remoteChain string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%d", localChain, remoteChain, time.Now().UnixNano()))))
	m.connections[id] = &ChainConnection{
		ID:          id,
		LocalChain:  localChain,
		RemoteChain: remoteChain,
		OpenedAt:    time.Now(),
	}
	return id
}

// CloseConnection marks a connection as closed.
func (m *ConnectionManager) CloseConnection(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.connections[id]
	if !ok {
		return fmt.Errorf("connection not found")
	}
	if c.Closed {
		return fmt.Errorf("connection already closed")
	}
	c.Closed = true
	c.ClosedAt = time.Now()
	return nil
}

// GetConnection retrieves connection details by ID.
func (m *ConnectionManager) GetConnection(id string) (*ChainConnection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.connections[id]
	return c, ok
}

// ListConnections returns all known connections.
func (m *ConnectionManager) ListConnections() []*ChainConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*ChainConnection, 0, len(m.connections))
	for _, c := range m.connections {
		out = append(out, c)
	}
	return out
}
