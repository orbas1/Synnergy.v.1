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
	Relayers    map[string]struct{}
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
func (m *ChainConnectionManager) Open(local, remote, relayer string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	id := m.nextID
	relayers := make(map[string]struct{})
	if relayer != "" {
		relayers[relayer] = struct{}{}
	}
	m.connections[id] = &ChainConnection{ID: id, LocalChain: local, RemoteChain: remote, Open: true, Relayers: relayers}
	return id
}

// Close terminates an existing connection. Only an authorized relayer may close the connection.
func (m *ChainConnectionManager) Close(id int, relayer string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.connections[id]
	if !ok {
		return errors.New("connection not found")
	}
	if _, authorized := c.Relayers[relayer]; !authorized {
		return errors.New("relayer not authorized")
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

// AuthorizeRelayer adds an address to the connection's relayer whitelist.
func (m *ChainConnectionManager) AuthorizeRelayer(id int, addr string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.connections[id]
	if !ok {
		return errors.New("connection not found")
	}
	if c.Relayers == nil {
		c.Relayers = make(map[string]struct{})
	}
	c.Relayers[addr] = struct{}{}
	return nil
}

// RevokeRelayer removes an address from the connection's whitelist.
func (m *ChainConnectionManager) RevokeRelayer(id int, addr string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.connections[id]
	if !ok {
		return errors.New("connection not found")
	}
	delete(c.Relayers, addr)
	return nil
}

// IsRelayerAuthorized checks whether the given address is authorized for the connection.
// It returns false if the connection does not exist or the relayer is not whitelisted.
func (m *ChainConnectionManager) IsRelayerAuthorized(id int, addr string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.connections[id]
	if !ok {
		return false
	}
	_, authorized := c.Relayers[addr]
	return authorized
}

// Remove deletes a connection from the manager.
// It returns an error if the connection cannot be found.
func (m *ChainConnectionManager) Remove(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.connections[id]; !ok {
		return errors.New("connection not found")
	}
	delete(m.connections, id)
	return nil
}
