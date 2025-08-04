package core

import "sync"

// InitService bootstraps a ledger using the replication subsystem. It is a
// thin wrapper around Replicator providing start and stop controls for the CLI.
type InitService struct {
	mu         sync.Mutex
	replicator *Replicator
	running    bool
}

// NewInitService creates a new initialization service bound to a replicator.
func NewInitService(r *Replicator) *InitService {
	return &InitService{replicator: r}
}

// Start bootstraps the ledger and starts replication.
func (i *InitService) Start() {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.running {
		return
	}
	i.replicator.Start()
	i.running = true
}

// Stop stops the initialization service and replication.
func (i *InitService) Stop() {
	i.mu.Lock()
	defer i.mu.Unlock()
	if !i.running {
		return
	}
	i.replicator.Stop()
	i.running = false
}
