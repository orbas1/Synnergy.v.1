package core

import (
	"errors"
	"net"
	"sync"
	"time"

	ilog "synnergy/internal/log"
)

// Connection represents a lightweight placeholder for an outbound connection.
type Connection struct {
	ID   string
	Conn net.Conn
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
func (p *ConnectionPool) Acquire(addr string) (*Connection, error) {
	p.mu.Lock()
	if c, ok := p.conns[addr]; ok {
		p.mu.Unlock()
		ilog.Info("conn_reuse", "id", addr)
		return c, nil
	}
	if len(p.conns) >= p.max {
		p.mu.Unlock()
		ilog.Error("conn_acquire", "error", "pool_exhausted")
		return nil, errors.New("connection pool exhausted")
	}
	p.mu.Unlock()

	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		ilog.Error("conn_dial_fail", "id", addr, "error", err)
		return nil, err
	}
	c := &Connection{ID: addr, Conn: conn}

	p.mu.Lock()
	p.conns[addr] = c
	p.mu.Unlock()

	ilog.Info("conn_new", "id", addr)
	return c, nil
}

// Release removes a connection from the pool.
func (p *ConnectionPool) Release(id string) {
	p.mu.Lock()
	c, ok := p.conns[id]
	if ok {
		delete(p.conns, id)
	}
	p.mu.Unlock()
	if ok && c.Conn != nil {
		_ = c.Conn.Close()
	}
	ilog.Info("conn_release", "id", id)
}

// Size returns the current number of active connections.
func (p *ConnectionPool) Size() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	size := len(p.conns)
	ilog.Info("conn_size", "size", size)
	return size
}

// Dial is a convenience wrapper around Acquire used by the CLI. It either
// returns an existing connection for the address or creates a new one if
// capacity allows.
func (p *ConnectionPool) Dial(addr string) (*Connection, error) {
	c, err := p.Acquire(addr)
	if err == nil {
		ilog.Info("conn_dial", "id", addr)
	}
	return c, err
}

// Close removes all connections from the pool, effectively resetting it.
func (p *ConnectionPool) Close() {
	p.mu.Lock()
	conns := p.conns
	p.conns = make(map[string]*Connection)
	p.mu.Unlock()
	for _, c := range conns {
		if c.Conn != nil {
			_ = c.Conn.Close()
		}
	}
	ilog.Info("conn_close")
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
	stats := PoolStats{Active: len(p.conns), Capacity: p.max}
	ilog.Info("conn_stats", "active", stats.Active, "capacity", stats.Capacity)
	return stats
}
