package watchtower

import (
	"context"
	"errors"
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
	JitterMS        float64   // Network jitter in milliseconds
	Downtime        time.Duration
	Alerts          []string
}

// ForkEvent captures details of a detected fork.
type ForkEvent struct {
	Height     uint64
	Hash       string
	Detected   time.Time
	Resolved   bool
	ResolvedAt time.Time
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

// ValidateMetrics verifies basic invariants for watchtower metrics.
func ValidateMetrics(m Metrics) error {
	if m.CPUUsage < 0 || m.CPUUsage > 100 {
		return errors.New("cpu usage must be between 0 and 100 percent")
	}
	if m.PeerCount < 0 {
		return errors.New("peer count cannot be negative")
	}
	return nil
}
