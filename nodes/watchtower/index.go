package watchtower

import (
	"context"
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
