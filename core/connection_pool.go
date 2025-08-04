package core

import (
	"errors"
	"sync"
)

// Connection represents a lightweight placeholder for an outbound connection.
type Connection struct {
	ID string
}

// ConnectionPool manages reusable connections to peers.
type ConnectionPool struct {
	mu    sync.Mutex
	conns map[string]*Connection
	max   int
}

// NewConnectionPool creates a pool with a maximum number of connections.
func NewConnectionPool(max int) *ConnectionPool {
	if max <= 0 {
		max = 1
	}
	return &ConnectionPool{conns: make(map[string]*Connection), max: max}
}

// Acquire returns an existing connection for the id or creates a new one if
// capacity allows.
func (p *ConnectionPool) Acquire(id string) (*Connection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if c, ok := p.conns[id]; ok {
		return c, nil
	}
	if len(p.conns) >= p.max {
		return nil, errors.New("connection pool exhausted")
	}
	c := &Connection{ID: id}
	p.conns[id] = c
	return c, nil
}

// Release removes a connection from the pool.
func (p *ConnectionPool) Release(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.conns, id)
}

// Size returns the current number of active connections.
func (p *ConnectionPool) Size() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.conns)
}
