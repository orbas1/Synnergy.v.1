package watchtower

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// BaseNode defines minimal behaviour required by all node implementations.
type BaseNode interface {
	// ID returns the unique identifier of the node.
	ID() string
}

// Metrics captures a snapshot of node health statistics.  Values are intended
// to be lightweight and easily serialisable for remote reporting.
type Metrics struct {
	CPUUsage        float64   // Percentage of CPU utilised by the node process
	MemoryUsage     uint64    // Bytes of RAM in use by the process
	PeerCount       int       // Number of peers currently connected
	LastBlockHeight uint64    // Height of the most recently observed block
	Timestamp       time.Time // Time the metrics were captured
}

// WatchtowerNode defines the operations exposed by a watchtower node.  These
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

// ErrStaleMetrics indicates that the supplied metrics snapshot is too old to be
// considered reliable by operators or automated tooling.  The caller should
// request a fresh sample from the watchtower before making decisions based on
// the metrics payload.
var ErrStaleMetrics = errors.New("watchtower: metrics snapshot is stale")

// ValidateMetrics performs lightweight sanity checking on the metrics gathered
// by a watchtower node.  The helper is intentionally opinionated so that CLI,
// web and automated operators all enforce the same health expectations without
// duplicating guard rails across projects.
//
// The validation rules favour defensive defaults: timestamps must be recent,
// utilisation percentages bounded and any negative counters rejected outright.
// This keeps downstream automation from acting on corrupt data when nodes are
// misconfigured or under active attack.
func ValidateMetrics(now time.Time, m Metrics) error {
	if m.Timestamp.IsZero() {
		return errors.New("watchtower: metrics timestamp missing")
	}
	if m.Timestamp.After(now.Add(30 * time.Second)) {
		return fmt.Errorf("watchtower: metrics timestamp %s is in the future", m.Timestamp)
	}
	age := now.Sub(m.Timestamp)
	if age < 0 {
		age = -age
	}
	if age > 5*time.Minute {
		return ErrStaleMetrics
	}
	if m.CPUUsage < 0 || m.CPUUsage > 100 {
		return fmt.Errorf("watchtower: cpu usage %.2f out of range", m.CPUUsage)
	}
	if m.MemoryUsage == 0 {
		return errors.New("watchtower: memory usage missing")
	}
	if m.PeerCount < 0 {
		return errors.New("watchtower: negative peer count")
	}
	if m.LastBlockHeight == 0 {
		return errors.New("watchtower: last block height missing")
	}
	return nil
}
