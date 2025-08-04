package core

import "sync"

// Replicator propagates blocks and snapshots to peers. This simplified
// implementation tracks replication state for CLI testing.
type Replicator struct {
	mu         sync.RWMutex
	ledger     *Ledger
	running    bool
	replicated map[string]bool
}

// NewReplicator constructs a Replicator bound to a ledger.
func NewReplicator(l *Ledger) *Replicator {
	return &Replicator{ledger: l, replicated: make(map[string]bool)}
}

// Start begins replication processes.
func (r *Replicator) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.running = true
}

// Stop halts replication.
func (r *Replicator) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.running = false
}

// Status reports whether the replicator is active.
func (r *Replicator) Status() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.running
}

// ReplicateBlock marks a block hash as replicated. In a full implementation
// this would gossip the block to peers.
func (r *Replicator) ReplicateBlock(hash string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.running {
		return false
	}
	r.replicated[hash] = true
	return true
}
