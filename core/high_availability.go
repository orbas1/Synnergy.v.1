package core

import (
	"sync"
	"time"
)

// FailoverManager tracks node heartbeats to provide high availability through
// automatic promotion of backup nodes when the primary becomes unresponsive.
type FailoverManager struct {
	mu      sync.RWMutex
	primary string
	nodes   map[string]time.Time // last heartbeat for each node
	timeout time.Duration
}

// NewFailoverManager creates a FailoverManager with a primary node identifier
// and a timeout indicating how long a node may miss heartbeats before being
// considered offline.
func NewFailoverManager(primary string, timeout time.Duration) *FailoverManager {
	return &FailoverManager{
		primary: primary,
		nodes:   map[string]time.Time{primary: time.Now()},
		timeout: timeout,
	}
}

// RegisterBackup adds a new backup node to the manager.
func (m *FailoverManager) RegisterBackup(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nodes[id] = time.Now()
}

// Heartbeat records a heartbeat for the specified node.
func (m *FailoverManager) Heartbeat(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nodes[id] = time.Now()
}

// RemoveNode removes a node from consideration. If the primary is removed the
// next call to Active will promote the freshest backup.
func (m *FailoverManager) RemoveNode(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.nodes, id)
	if id == m.primary {
		m.primary = ""
	}
}

// Active returns the identifier of the node currently acting as primary.  If the
// existing primary has not sent a heartbeat within the timeout, the most recent
// backup node is promoted.
func (m *FailoverManager) Active() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if hb, ok := m.nodes[m.primary]; ok {
		if time.Since(hb) <= m.timeout {
			return m.primary
		}
	}

	var candidate string
	var latest time.Time
	for id, hb := range m.nodes {
		if id == m.primary {
			continue
		}
		if hb.After(latest) {
			candidate = id
			latest = hb
		}
	}
	if candidate != "" {
		m.primary = candidate
	}
	return m.primary
}
