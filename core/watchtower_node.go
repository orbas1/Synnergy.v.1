package core

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"synnergy/nodesextra/watchtower"
)

// Watchtower implements the watchtower.WatchtowerNode interface.  It observes
// the network, records system health metrics and reports detected forks.
type Watchtower struct {
	id       string
	firewall *Firewall
	health   *SystemHealthLogger
	logger   *log.Logger

	mu      sync.RWMutex
	running bool
	cancel  context.CancelFunc
}

// NewWatchtowerNode constructs a Watchtower with the provided identifier.
func NewWatchtowerNode(id string, logger *log.Logger) *Watchtower {
	return &Watchtower{
		id:       id,
		firewall: NewFirewall(),
		health:   NewSystemHealthLogger(),
		logger:   logger,
	}
}

// ID returns the unique identifier of the node.
func (w *Watchtower) ID() string { return w.id }

// Start begins monitoring routines for the watchtower node.
func (w *Watchtower) Start(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.running {
		return errors.New("watchtower already running")
	}
	var c context.Context
	c, w.cancel = context.WithCancel(ctx)
	w.running = true
	go w.monitorLoop(c)
	return nil
}

// monitorLoop periodically collects system health metrics.
func (w *Watchtower) monitorLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m := w.health.Collect(0, 0)
			if w.logger != nil {
				w.logger.Printf("watchtower metrics: %+v", m)
			}
		case <-ctx.Done():
			return
		}
	}
}

// Stop halts monitoring routines for the watchtower node.
func (w *Watchtower) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.running {
		return nil
	}
	w.cancel()
	w.running = false
	return nil
}

// ReportFork records details of a detected fork event.
func (w *Watchtower) ReportFork(height uint64, hash string) {
	if w.logger != nil {
		w.logger.Printf("fork detected at height %d hash %s", height, hash)
	}
}

// Metrics returns the latest snapshot of system health data.
func (w *Watchtower) Metrics() watchtower.Metrics {
	return w.health.Snapshot()
}

// Firewall exposes the internal firewall instance for rule management.
func (w *Watchtower) Firewall() *Firewall { return w.firewall }
