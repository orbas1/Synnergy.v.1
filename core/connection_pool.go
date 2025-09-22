package core

import (
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	ilog "synnergy/internal/log"
)

// Connection represents a lightweight placeholder for an outbound connection.
type Connection struct {
	ID   string
	Conn net.Conn
}

// PoolOptions configure connection reuse and health monitoring behaviour.
type PoolOptions struct {
	Max                 int
	DialTimeout         time.Duration
	IdleTimeout         time.Duration
	HealthCheckInterval time.Duration
	TLSConfig           *tls.Config
}

type pooledConn struct {
	conn      *Connection
	lastUsed  time.Time
	unhealthy bool
}

type poolMetrics struct {
	created    atomic.Uint64
	reused     atomic.Uint64
	dialFailed atomic.Uint64
	closedIdle atomic.Uint64
}

// ConnectionPool manages reusable connections to peers.
type ConnectionPool struct {
	mu      sync.Mutex
	conns   map[string]*pooledConn
	max     int
	opts    PoolOptions
	quit    chan struct{}
	wg      sync.WaitGroup
	once    sync.Once
	metrics poolMetrics
}

// NewConnectionPool creates a pool with a maximum number of connections.
func NewConnectionPool(max int) *ConnectionPool {
	return NewConnectionPoolWithOptions(PoolOptions{Max: max})
}

// NewConnectionPoolWithOptions creates a pool using the supplied options.
func NewConnectionPoolWithOptions(opts PoolOptions) *ConnectionPool {
	if opts.Max <= 0 {
		opts.Max = 1
	}
	if opts.DialTimeout <= 0 {
		opts.DialTimeout = 3 * time.Second
	}
	if opts.IdleTimeout < 0 {
		opts.IdleTimeout = 0
	}
	if opts.HealthCheckInterval <= 0 {
		if opts.IdleTimeout > 0 {
			opts.HealthCheckInterval = opts.IdleTimeout / 2
		}
		if opts.HealthCheckInterval <= 0 {
			opts.HealthCheckInterval = time.Minute
		}
	}
	pool := &ConnectionPool{
		conns: make(map[string]*pooledConn),
		max:   opts.Max,
		opts:  opts,
		quit:  make(chan struct{}),
	}
	pool.wg.Add(1)
	go pool.healthLoop()
	return pool
}

// Acquire returns an existing connection for the id or creates a new one if
// capacity allows. Unhealthy or stale connections are closed before creating a
// replacement.
func (p *ConnectionPool) Acquire(addr string) (*Connection, error) {
	now := time.Now()
	p.mu.Lock()
	if pc, ok := p.conns[addr]; ok {
		if pc.unhealthy || (p.opts.IdleTimeout > 0 && now.Sub(pc.lastUsed) > p.opts.IdleTimeout) {
			delete(p.conns, addr)
			p.mu.Unlock()
			if pc.conn != nil && pc.conn.Conn != nil {
				_ = pc.conn.Conn.Close()
			}
			p.metrics.closedIdle.Add(1)
			return p.dial(addr, now)
		}
		pc.lastUsed = now
		p.mu.Unlock()
		p.metrics.reused.Add(1)
		ilog.Info("conn_reuse", "id", addr)
		return pc.conn, nil
	}
	if len(p.conns) >= p.max {
		p.mu.Unlock()
		p.metrics.dialFailed.Add(1)
		ilog.Error("conn_acquire", "id", addr, "error", "pool_exhausted")
		return nil, errors.New("connection pool exhausted")
	}
	p.mu.Unlock()
	return p.dial(addr, now)
}

func (p *ConnectionPool) dial(addr string, ts time.Time) (*Connection, error) {
	dialer := &net.Dialer{Timeout: p.opts.DialTimeout, KeepAlive: time.Second}
	var (
		raw net.Conn
		err error
	)
	if p.opts.TLSConfig != nil {
		raw, err = tls.DialWithDialer(dialer, "tcp", addr, p.opts.TLSConfig)
	} else {
		raw, err = dialer.Dial("tcp", addr)
	}
	if err != nil {
		p.metrics.dialFailed.Add(1)
		ilog.Error("conn_dial_fail", "id", addr, "error", err)
		return nil, err
	}
	conn := &Connection{ID: addr, Conn: raw}
	pc := &pooledConn{conn: conn, lastUsed: ts}

	p.mu.Lock()
	if len(p.conns) >= p.max {
		p.mu.Unlock()
		_ = raw.Close()
		p.metrics.dialFailed.Add(1)
		return nil, errors.New("connection pool exhausted")
	}
	p.conns[addr] = pc
	p.mu.Unlock()

	p.metrics.created.Add(1)
	ilog.Info("conn_new", "id", addr)
	return conn, nil
}

// Release removes a connection from the pool.
func (p *ConnectionPool) Release(id string) {
	p.mu.Lock()
	pc, ok := p.conns[id]
	if ok {
		delete(p.conns, id)
	}
	p.mu.Unlock()
	if ok && pc != nil && pc.conn != nil && pc.conn.Conn != nil {
		_ = pc.conn.Conn.Close()
	}
	ilog.Info("conn_release", "id", id)
}

// ReportFailure marks the connection as unhealthy so the next Acquire refreshes
// it.
func (p *ConnectionPool) ReportFailure(id string, err error) {
	p.mu.Lock()
	if pc, ok := p.conns[id]; ok {
		pc.unhealthy = true
	}
	p.mu.Unlock()
	if err != nil {
		ilog.Error("conn_failure", "id", id, "error", err)
	} else {
		ilog.Error("conn_failure", "id", id)
	}
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
	p.once.Do(func() {
		close(p.quit)
		p.wg.Wait()
		p.mu.Lock()
		conns := p.conns
		p.conns = make(map[string]*pooledConn)
		p.mu.Unlock()
		for _, pc := range conns {
			if pc != nil && pc.conn != nil && pc.conn.Conn != nil {
				_ = pc.conn.Conn.Close()
			}
		}
		ilog.Info("conn_close")
	})
}

func (p *ConnectionPool) healthLoop() {
	defer p.wg.Done()
	ticker := time.NewTicker(p.opts.HealthCheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			p.pruneStale(time.Now())
		case <-p.quit:
			return
		}
	}
}

func (p *ConnectionPool) pruneStale(now time.Time) {
	if p.opts.IdleTimeout <= 0 {
		return
	}
	p.mu.Lock()
	for addr, pc := range p.conns {
		if pc == nil {
			delete(p.conns, addr)
			continue
		}
		if pc.unhealthy || now.Sub(pc.lastUsed) > p.opts.IdleTimeout {
			delete(p.conns, addr)
			go func(c *Connection) {
				if c != nil && c.Conn != nil {
					_ = c.Conn.Close()
				}
			}(pc.conn)
			p.metrics.closedIdle.Add(1)
		}
	}
	p.mu.Unlock()
}

// PoolStats summarises connection usage and capacity for diagnostic output.
type PoolStats struct {
	Active       int
	Capacity     int
	Created      uint64
	Reused       uint64
	DialFailures uint64
	ClosedIdle   uint64
}

// Stats returns a snapshot of the pool's current usage.
func (p *ConnectionPool) Stats() PoolStats {
	p.mu.Lock()
	active := len(p.conns)
	capacity := p.max
	p.mu.Unlock()
	stats := PoolStats{
		Active:       active,
		Capacity:     capacity,
		Created:      p.metrics.created.Load(),
		Reused:       p.metrics.reused.Load(),
		DialFailures: p.metrics.dialFailed.Load(),
		ClosedIdle:   p.metrics.closedIdle.Load(),
	}
	ilog.Info("conn_stats",
		"active", stats.Active,
		"capacity", stats.Capacity,
		"created", stats.Created,
		"reused", stats.Reused,
		"dial_failures", stats.DialFailures,
		"closed_idle", stats.ClosedIdle,
	)
	return stats
}
