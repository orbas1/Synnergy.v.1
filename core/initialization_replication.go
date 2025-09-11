package core

import (
	"errors"
	"sync"
)

// InitService bootstraps a ledger using the replication subsystem. It is a
// thin wrapper around Replicator providing start and stop controls for the CLI.
type InitService struct {
	mu         sync.RWMutex
	replicator *Replicator
	running    bool
}

// NewInitService creates a new initialization service bound to a replicator.
func NewInitService(r *Replicator) *InitService {
	return &InitService{replicator: r}
}

var (
	ErrInitReplicatorNil  = errors.New("replicator is nil")
	ErrInitServiceRunning = errors.New("initialization already running")
	ErrInitServiceStopped = errors.New("initialization not running")
)

// Start bootstraps the ledger and starts replication.
func (i *InitService) Start() error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.replicator == nil {
		return ErrInitReplicatorNil
	}
	if i.running {
		return ErrInitServiceRunning
	}
	i.replicator.Start()
	i.running = true
	return nil
}

// Stop stops the initialization service and replication.
func (i *InitService) Stop() error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if !i.running {
		return ErrInitServiceStopped
	}
	i.replicator.Stop()
	i.running = false
	return nil
}

// Status reports whether the initialization service is running.
func (i *InitService) Status() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.running
}
