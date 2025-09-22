package synnergy

import (
	"sync"
	"time"
)

// FailoverManager tracks node heartbeats to provide high availability through
// automatic promotion of backup nodes when the primary becomes unresponsive.
type FailoverManager struct {
	mu         sync.RWMutex
	primary    string
	nodes      map[string]time.Time // last heartbeat for each node
	timeout    time.Duration
	failovers  int
	lastSwitch time.Time
}

// NewFailoverManager creates a FailoverManager with a primary node identifier
// and a timeout indicating how long a node may miss heartbeats before being
// considered offline.
func NewFailoverManager(primary string, timeout time.Duration) *FailoverManager {
	return &FailoverManager{
		primary: primary,
		nodes:   map[string]time.Time{primary: time.Now()},
		timeout: timeout,
		// Treat the initial primary assignment as a switch so callers
		// always receive a non-zero timestamp when inspecting
		// orchestration state.
		lastSwitch: time.Now(),
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

// Active returns the identifier of the node currently acting as primary. If the
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
		if candidate != m.primary {
			m.failovers++
			m.lastSwitch = time.Now()
		}
		m.primary = candidate
	}
	return m.primary
}

// FailoverSnapshot captures an immutable view of the failover manager's
// orchestration state for evaluation by other components such as the consensus
// hopper.
type FailoverSnapshot struct {
	Active     string
	Healthy    bool
	Failovers  int
	LastSwitch time.Time
	Timeout    time.Duration
}

// Snapshot returns the current orchestration snapshot including primary health
// and failover statistics. Callers should prefer Snapshot over peeking into the
// struct fields directly to avoid lock ordering issues.
func (m *FailoverManager) Snapshot() FailoverSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hb, ok := m.nodes[m.primary]
	healthy := false
	if ok {
		healthy = time.Since(hb) <= m.timeout
	}

	return FailoverSnapshot{
		Active:     m.primary,
		Healthy:    healthy,
		Failovers:  m.failovers,
		LastSwitch: m.lastSwitch,
		Timeout:    m.timeout,
	}
}
