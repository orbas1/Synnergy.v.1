package watchtower

import (
	"context"
	"errors"
	"sync"
	"time"
)

// BaseNode defines minimal behaviour required by all node implementations.
type BaseNode interface {
	// ID returns the unique identifier of the node.
	ID() string
}

// Metrics captures a snapshot of node health statistics. Values are intended
// to be lightweight and easily serialisable for remote reporting.
type Metrics struct {
	CPUUsage        float64   // Percentage of CPU utilised by the node process
	MemoryUsage     uint64    // Bytes of RAM in use by the process
	PeerCount       int       // Number of peers currently connected
	LastBlockHeight uint64    // Height of the most recently observed block
	Timestamp       time.Time // Time the metrics were captured
}

// WatchtowerNode defines the operations exposed by a watchtower node. These
// nodes observe the network, report forks and make system health metrics
// available to operators.
type WatchtowerNode interface {
	BaseNode

	// Start begins monitoring routines for the node.
	Start(ctx context.Context) error

	// Stop gracefully shuts down monitoring routines.
	Stop() error

	// ReportFork records details of a detected fork for later analysis.
	ReportFork(height uint64, hash string)

	// Metrics returns the latest snapshot of system health data.
	Metrics() Metrics
}

// BasicWatchtower provides an in-memory implementation of the WatchtowerNode
// interface.  It records forks and metrics for later inspection.
type BasicWatchtower struct {
	id      string
	mu      sync.RWMutex
	running bool
	metrics Metrics
	forks   []struct {
		Height uint64
		Hash   string
	}
}

// NewBasicWatchtower creates a new watchtower with the supplied identifier.
func NewBasicWatchtower(id string) *BasicWatchtower {
	return &BasicWatchtower{id: id}
}

// ID returns the identifier of the watchtower node.
func (w *BasicWatchtower) ID() string { return w.id }

// Start marks the node as running and initialises the metrics timestamp.
func (w *BasicWatchtower) Start(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.running {
		return errors.New("watchtower already running")
	}
	w.running = true
	w.metrics.Timestamp = time.Now()
	return nil
}

// Stop halts the node.
func (w *BasicWatchtower) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.running {
		return errors.New("watchtower not running")
	}
	w.running = false
	return nil
}

// ReportFork records a fork event.
func (w *BasicWatchtower) ReportFork(height uint64, hash string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.forks = append(w.forks, struct {
		Height uint64
		Hash   string
	}{height, hash})
}

// Metrics returns a snapshot of the latest metrics.
func (w *BasicWatchtower) Metrics() Metrics {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.metrics
}
