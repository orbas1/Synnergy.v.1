package core

import (
	"runtime"
	"sync"
	"time"

	"synnergy/Nodes/watchtower"
)

// SystemHealthLogger collects runtime metrics for a node and exposes snapshots
// for external consumption.  It is safe for concurrent use by multiple goroutines.
type SystemHealthLogger struct {
	mu   sync.RWMutex
	last watchtower.Metrics
}

// NewSystemHealthLogger returns an initialised SystemHealthLogger instance.
func NewSystemHealthLogger() *SystemHealthLogger {
	return &SystemHealthLogger{}
}

// Collect gathers metrics from the runtime and records them as the latest snapshot.
// Caller may supply optional peer count and block height information.
func (l *SystemHealthLogger) Collect(peerCount int, height uint64) watchtower.Metrics {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	m := watchtower.Metrics{
		CPUUsage:        float64(runtime.NumGoroutine()),
		MemoryUsage:     ms.Alloc,
		PeerCount:       peerCount,
		LastBlockHeight: height,
		Timestamp:       time.Now(),
	}

	l.mu.Lock()
	l.last = m
	l.mu.Unlock()
	return m
}

// Snapshot returns the most recently recorded metrics.
func (l *SystemHealthLogger) Snapshot() watchtower.Metrics {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.last
}
