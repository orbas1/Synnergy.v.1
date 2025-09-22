package core

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Sentinel errors exposed for CLI validation and unit tests.
var (
	ErrSandboxExists   = errors.New("sandbox already exists")
	ErrSandboxNotFound = errors.New("sandbox not found")
)

// SandboxEventType enumerates lifecycle transitions emitted by the sandbox
// manager. Events are consumed by the enterprise CLI, web telemetry panels and
// integration tests to provide real-time visibility of execution sandboxes.
type SandboxEventType string

const (
	SandboxEventStarted   SandboxEventType = "started"
	SandboxEventStopped   SandboxEventType = "stopped"
	SandboxEventReset     SandboxEventType = "reset"
	SandboxEventDeleted   SandboxEventType = "deleted"
	SandboxEventPurged    SandboxEventType = "purged"
	SandboxEventFailure   SandboxEventType = "failure"
	SandboxEventRestarted SandboxEventType = "restarted"
	SandboxEventHeartbeat SandboxEventType = "heartbeat"
)

// SandboxInfo holds runtime limits and state for a single sandboxed contract
// execution environment.
type SandboxInfo struct {
	ID            string
	ContractAddr  string
	GasLimit      uint64
	MemoryLimit   uint64
	Active        bool
	LastReset     time.Time
	CreatedAt     time.Time
	LastHeartbeat time.Time
	FailureCount  uint64
	RestartCount  uint64
	LastError     string
}

// SandboxEvent captures a lifecycle notification for subscribers.
type SandboxEvent struct {
	Type      SandboxEventType
	Info      SandboxInfo
	Timestamp time.Time
	Err       error
}

// SandboxWatcher receives sandbox events.
type SandboxWatcher func(SandboxEvent)

// SandboxMetrics exposes aggregate counters for observability dashboards.
type SandboxMetrics struct {
	TotalCreated uint64
	Active       int64
	Failures     uint64
	Restarts     uint64
	Purged       uint64
}

type sandboxMetrics struct {
	totalCreated uint64
	active       int64
	failures     uint64
	restarts     uint64
	purged       uint64
}

// SandboxManager manages multiple sandboxes used to isolate contract
// execution. It is safe for concurrent use.
type SandboxManager struct {
	mu        sync.RWMutex
	sandboxes map[string]*SandboxInfo
	ttl       time.Duration
	watchers  []SandboxWatcher
	metrics   sandboxMetrics
}

// NewSandboxManager returns an empty manager. A default time-to-live of one
// hour is applied to inactive sandboxes unless overridden. Expired sandboxes
// are removed when PurgeInactive is invoked.
func NewSandboxManager(ttls ...time.Duration) *SandboxManager {
	ttl := time.Hour
	if len(ttls) > 0 && ttls[0] > 0 {
		ttl = ttls[0]
	}
	return &SandboxManager{sandboxes: make(map[string]*SandboxInfo), ttl: ttl}
}

// RegisterWatcher attaches a sandbox watcher. Callers should register
// watchers before invoking other manager methods.
func (m *SandboxManager) RegisterWatcher(w SandboxWatcher) {
	if w == nil {
		return
	}
	m.mu.Lock()
	m.watchers = append(m.watchers, w)
	m.mu.Unlock()
}

func (m *SandboxManager) snapshotWatchers() []SandboxWatcher {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.watchers) == 0 {
		return nil
	}
	out := make([]SandboxWatcher, len(m.watchers))
	copy(out, m.watchers)
	return out
}

func (m *SandboxManager) emit(event SandboxEvent) {
	for _, w := range m.snapshotWatchers() {
		func() {
			defer func() { _ = recover() }()
			w(event)
		}()
	}
}

// StartSandbox creates a new sandbox for the given contract address.
func (m *SandboxManager) StartSandbox(id, contractAddr string, gasLimit, memoryLimit uint64) (*SandboxInfo, error) {
	now := time.Now()
	m.mu.Lock()
	if _, exists := m.sandboxes[id]; exists {
		m.mu.Unlock()
		return nil, ErrSandboxExists
	}
	sb := &SandboxInfo{
		ID:            id,
		ContractAddr:  contractAddr,
		GasLimit:      gasLimit,
		MemoryLimit:   memoryLimit,
		Active:        true,
		CreatedAt:     now,
		LastReset:     now,
		LastHeartbeat: now,
	}
	m.sandboxes[id] = sb
	m.mu.Unlock()

	atomic.AddUint64(&m.metrics.totalCreated, 1)
	atomic.AddInt64(&m.metrics.active, 1)
	m.emit(SandboxEvent{Type: SandboxEventStarted, Info: *sb, Timestamp: now})
	return sb, nil
}

// StopSandbox deactivates a sandbox.
func (m *SandboxManager) StopSandbox(id string) error {
	now := time.Now()
	m.mu.Lock()
	sb, ok := m.sandboxes[id]
	if !ok {
		m.mu.Unlock()
		return ErrSandboxNotFound
	}
	wasActive := sb.Active
	sb.Active = false
	sb.LastHeartbeat = now
	snapshot := *sb
	m.mu.Unlock()

	if wasActive {
		atomic.AddInt64(&m.metrics.active, -1)
	}
	m.emit(SandboxEvent{Type: SandboxEventStopped, Info: snapshot, Timestamp: now})
	return nil
}

// ResetSandbox updates the LastReset timestamp for a sandbox.
func (m *SandboxManager) ResetSandbox(id string) error {
	now := time.Now()
	m.mu.Lock()
	sb, ok := m.sandboxes[id]
	if !ok {
		m.mu.Unlock()
		return ErrSandboxNotFound
	}
	sb.LastReset = now
	sb.LastHeartbeat = now
	snapshot := *sb
	m.mu.Unlock()

	m.emit(SandboxEvent{Type: SandboxEventReset, Info: snapshot, Timestamp: now})
	return nil
}

// DeleteSandbox removes a sandbox entirely, releasing its resources.
func (m *SandboxManager) DeleteSandbox(id string) error {
	m.mu.Lock()
	sb, ok := m.sandboxes[id]
	if !ok {
		m.mu.Unlock()
		return ErrSandboxNotFound
	}
	wasActive := sb.Active
	snapshot := *sb
	delete(m.sandboxes, id)
	m.mu.Unlock()

	if wasActive {
		atomic.AddInt64(&m.metrics.active, -1)
	}
	m.emit(SandboxEvent{Type: SandboxEventDeleted, Info: snapshot, Timestamp: time.Now()})
	return nil
}

// RecordFailure marks a sandbox as failed and stores the error for later
// inspection. The sandbox is deactivated until RestartSandbox is invoked.
func (m *SandboxManager) RecordFailure(id string, failure error) error {
	now := time.Now()
	m.mu.Lock()
	sb, ok := m.sandboxes[id]
	if !ok {
		m.mu.Unlock()
		return ErrSandboxNotFound
	}
	if sb.Active {
		atomic.AddInt64(&m.metrics.active, -1)
	}
	sb.Active = false
	sb.LastHeartbeat = now
	sb.FailureCount++
	if failure != nil {
		sb.LastError = failure.Error()
	} else {
		sb.LastError = ""
	}
	snapshot := *sb
	m.mu.Unlock()

	atomic.AddUint64(&m.metrics.failures, 1)
	m.emit(SandboxEvent{Type: SandboxEventFailure, Info: snapshot, Timestamp: now, Err: failure})
	return nil
}

// RestartSandbox reactivates a sandbox following a failure or manual stop.
func (m *SandboxManager) RestartSandbox(id string) error {
	now := time.Now()
	m.mu.Lock()
	sb, ok := m.sandboxes[id]
	if !ok {
		m.mu.Unlock()
		return ErrSandboxNotFound
	}
	if !sb.Active {
		atomic.AddInt64(&m.metrics.active, 1)
	}
	sb.Active = true
	sb.RestartCount++
	sb.LastReset = now
	sb.LastHeartbeat = now
	snapshot := *sb
	m.mu.Unlock()

	atomic.AddUint64(&m.metrics.restarts, 1)
	m.emit(SandboxEvent{Type: SandboxEventRestarted, Info: snapshot, Timestamp: now})
	return nil
}

// Heartbeat updates the liveness timestamp for the sandbox. It is typically
// invoked by the VM after a successful execution round.
func (m *SandboxManager) Heartbeat(id string) error {
	now := time.Now()
	m.mu.Lock()
	sb, ok := m.sandboxes[id]
	if !ok {
		m.mu.Unlock()
		return ErrSandboxNotFound
	}
	sb.LastHeartbeat = now
	snapshot := *sb
	m.mu.Unlock()

	m.emit(SandboxEvent{Type: SandboxEventHeartbeat, Info: snapshot, Timestamp: now})
	return nil
}

// PurgeInactive deletes sandboxes that have been stopped and remained idle
// beyond the configured TTL.  It is safe to call concurrently with other
// manager operations.
func (m *SandboxManager) PurgeInactive() {
	now := time.Now()
	var purged []SandboxInfo

	m.mu.Lock()
	for id, sb := range m.sandboxes {
		if sb.Active {
			continue
		}
		if now.Sub(sb.LastHeartbeat) > m.ttl {
			purged = append(purged, *sb)
			delete(m.sandboxes, id)
		}
	}
	m.mu.Unlock()

	if len(purged) == 0 {
		return
	}
	atomic.AddUint64(&m.metrics.purged, uint64(len(purged)))
	for _, sb := range purged {
		m.emit(SandboxEvent{Type: SandboxEventPurged, Info: sb, Timestamp: now})
	}
}

// SandboxStatus returns sandbox information by ID.
func (m *SandboxManager) SandboxStatus(id string) (*SandboxInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sb, ok := m.sandboxes[id]
	return sb, ok
}

// ListSandboxes returns all sandboxes managed by this instance.
func (m *SandboxManager) ListSandboxes() []*SandboxInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*SandboxInfo, 0, len(m.sandboxes))
	for _, sb := range m.sandboxes {
		out = append(out, sb)
	}
	return out
}

// Metrics returns a snapshot of the sandbox counters.
func (m *SandboxManager) Metrics() SandboxMetrics {
	return SandboxMetrics{
		TotalCreated: atomic.LoadUint64(&m.metrics.totalCreated),
		Active:       atomic.LoadInt64(&m.metrics.active),
		Failures:     atomic.LoadUint64(&m.metrics.failures),
		Restarts:     atomic.LoadUint64(&m.metrics.restarts),
		Purged:       atomic.LoadUint64(&m.metrics.purged),
	}
}
