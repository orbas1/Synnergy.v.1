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

// Dial is a convenience wrapper around Acquire used by the CLI. It either
// returns an existing connection for the address or creates a new one if
// capacity allows.
func (p *ConnectionPool) Dial(addr string) (*Connection, error) {
	return p.Acquire(addr)
}

// Close removes all connections from the pool, effectively resetting it.
func (p *ConnectionPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.conns = make(map[string]*Connection)
}

// PoolStats summarises connection usage and capacity for diagnostic output.
type PoolStats struct {
	Active   int
	Capacity int
}

// Stats returns a snapshot of the pool's current usage.
func (p *ConnectionPool) Stats() PoolStats {
	p.mu.Lock()
	defer p.mu.Unlock()
	return PoolStats{Active: len(p.conns), Capacity: p.max}
}
