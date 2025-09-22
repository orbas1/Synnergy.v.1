package core

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"synnergy/internal/nodes"
)

const replicationQueueSize = 128

type replicationRequest struct {
	hash    string
	attempt int
}

// ReplicationRecord captures acknowledgement progress for a block hash.
type ReplicationRecord struct {
	Hash         string
	Attempts     int
	LastAttempt  time.Time
	Acknowledged map[nodes.Address]time.Time
	Completed    bool
}

type replicationMetrics struct {
	enqueued     atomic.Uint64
	acknowledged atomic.Uint64
	retries      atomic.Uint64
	dropped      atomic.Uint64
}

// ReplicationMetrics exposes runtime counters for observability.
type ReplicationMetrics struct {
	Enqueued     uint64
	Acknowledged uint64
	Retries      uint64
	Dropped      uint64
	Pending      int
}

// Replicator propagates blocks and snapshots to peers. It queues replication
// requests, tracks acknowledgements, and automatically retries until all peers
// confirm receipt or the retry budget is exhausted.
type Replicator struct {
	mu         sync.RWMutex
	ledger     *Ledger
	running    bool
	peers      map[nodes.Address]struct{}
	records    map[string]*ReplicationRecord
	queue      chan replicationRequest
	quit       chan struct{}
	wg         sync.WaitGroup
	retryLimit int
	retryDelay time.Duration
	timeout    time.Duration
	metrics    replicationMetrics
}

// NewReplicator constructs a Replicator bound to a ledger.
func NewReplicator(l *Ledger) *Replicator {
	return &Replicator{
		ledger:     l,
		peers:      make(map[nodes.Address]struct{}),
		records:    make(map[string]*ReplicationRecord),
		retryLimit: 3,
		retryDelay: 200 * time.Millisecond,
		timeout:    500 * time.Millisecond,
	}
}

// Start begins replication processes.
func (r *Replicator) Start() {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return
	}
	r.queue = make(chan replicationRequest, replicationQueueSize)
	r.quit = make(chan struct{})
	r.running = true
	r.wg.Add(1)
	go r.run()
	r.mu.Unlock()
}

// Stop halts replication.
func (r *Replicator) Stop() {
	r.mu.Lock()
	if !r.running {
		r.mu.Unlock()
		return
	}
	quit := r.quit
	r.running = false
	r.quit = nil
	r.mu.Unlock()
	close(quit)
	r.wg.Wait()
	r.mu.Lock()
	r.queue = nil
	r.mu.Unlock()
}

// Status reports whether the replicator is active.
func (r *Replicator) Status() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.running
}

// RegisterPeer adds a peer to the acknowledgement set.
func (r *Replicator) RegisterPeer(addr nodes.Address) {
	if addr == "" {
		return
	}
	r.mu.Lock()
	if _, ok := r.peers[addr]; ok {
		r.mu.Unlock()
		return
	}
	r.peers[addr] = struct{}{}
	hashes := make([]string, 0, len(r.records))
	for hash, rec := range r.records {
		if rec.Acknowledged == nil {
			rec.Acknowledged = make(map[nodes.Address]time.Time)
		}
		rec.Acknowledged[addr] = time.Time{}
		rec.Completed = false
		hashes = append(hashes, hash)
	}
	running := r.running
	r.mu.Unlock()
	if running {
		for _, hash := range hashes {
			_ = r.tryEnqueue(replicationRequest{hash: hash})
		}
	}
}

// RemovePeer removes acknowledgement tracking for the provided address.
func (r *Replicator) RemovePeer(addr nodes.Address) {
	r.mu.Lock()
	delete(r.peers, addr)
	for _, rec := range r.records {
		delete(rec.Acknowledged, addr)
	}
	r.mu.Unlock()
}

// ReplicateBlock schedules replication for the provided block hash. It returns
// false if the replicator is stopped or the block is unknown.
func (r *Replicator) ReplicateBlock(hash string) bool {
	if hash == "" {
		return false
	}
	if r.ledger != nil && !r.ledger.HasBlock(hash) {
		return false
	}
	if err := r.tryEnqueue(replicationRequest{hash: hash}); err != nil {
		return false
	}
	return true
}

var errReplicatorStopped = errors.New("replicator stopped")

func (r *Replicator) tryEnqueue(req replicationRequest) error {
	r.mu.RLock()
	running := r.running
	queue := r.queue
	quit := r.quit
	timeout := r.timeout
	r.mu.RUnlock()
	if !running || queue == nil {
		return errReplicatorStopped
	}
	if timeout <= 0 {
		select {
		case queue <- req:
			r.metrics.enqueued.Add(1)
			return nil
		case <-quit:
			return errReplicatorStopped
		}
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case queue <- req:
		r.metrics.enqueued.Add(1)
		return nil
	case <-timer.C:
		r.metrics.dropped.Add(1)
		return errors.New("enqueue timeout")
	case <-quit:
		return errReplicatorStopped
	}
}

func (r *Replicator) run() {
	defer r.wg.Done()
	for {
		select {
		case req := <-r.queue:
			r.handleRequest(req)
		case <-r.quit:
			return
		}
	}
}

func (r *Replicator) handleRequest(req replicationRequest) {
	now := time.Now()
	r.mu.Lock()
	rec := r.ensureRecordLocked(req.hash)
	rec.Attempts = req.attempt + 1
	rec.LastAttempt = now
	pending := r.pendingPeersLocked(rec)
	r.mu.Unlock()
	if len(pending) == 0 {
		return
	}
	if req.attempt >= r.retryLimit {
		return
	}
	r.scheduleRetry(replicationRequest{hash: req.hash, attempt: req.attempt + 1})
}

func (r *Replicator) ensureRecordLocked(hash string) *ReplicationRecord {
	rec := r.records[hash]
	if rec == nil {
		rec = &ReplicationRecord{
			Hash:         hash,
			Acknowledged: make(map[nodes.Address]time.Time),
		}
		for peer := range r.peers {
			rec.Acknowledged[peer] = time.Time{}
		}
		if len(r.peers) == 0 {
			rec.Completed = true
		}
		r.records[hash] = rec
	} else {
		for peer := range r.peers {
			if _, ok := rec.Acknowledged[peer]; !ok {
				rec.Acknowledged[peer] = time.Time{}
				rec.Completed = false
			}
		}
	}
	return rec
}

func (r *Replicator) pendingPeersLocked(rec *ReplicationRecord) []nodes.Address {
	pending := make([]nodes.Address, 0)
	for peer, ts := range rec.Acknowledged {
		if ts.IsZero() {
			pending = append(pending, peer)
		}
	}
	if len(pending) == 0 && len(rec.Acknowledged) > 0 {
		rec.Completed = true
	}
	return pending
}

func (r *Replicator) scheduleRetry(req replicationRequest) {
	r.metrics.retries.Add(1)
	delay := r.retryDelay
	if delay <= 0 {
		delay = 100 * time.Millisecond
	}
	delay = delay * time.Duration(1<<req.attempt)
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-timer.C:
			if err := r.tryEnqueue(req); err != nil {
				r.metrics.dropped.Add(1)
			}
		case <-r.quit:
			return
		}
	}()
}

// MarkAcknowledged marks a peer as having received the replication for the hash.
func (r *Replicator) MarkAcknowledged(hash string, addr nodes.Address) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec := r.records[hash]
	if rec == nil {
		return false
	}
	if _, ok := rec.Acknowledged[addr]; !ok {
		return false
	}
	if rec.Acknowledged[addr].IsZero() {
		rec.Acknowledged[addr] = time.Now()
		if len(r.pendingPeersLocked(rec)) == 0 {
			rec.Completed = true
		}
		r.metrics.acknowledged.Add(1)
	}
	return true
}

// Replicated reports whether a block has been replicated to all known peers.
func (r *Replicator) Replicated(hash string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec := r.records[hash]
	return rec != nil && rec.Completed
}

// ReplicationStatus returns a copy of the replication record for the hash.
func (r *Replicator) ReplicationStatus(hash string) (ReplicationRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec := r.records[hash]
	if rec == nil {
		return ReplicationRecord{}, false
	}
	cp := ReplicationRecord{
		Hash:        rec.Hash,
		Attempts:    rec.Attempts,
		LastAttempt: rec.LastAttempt,
		Completed:   rec.Completed,
	}
	cp.Acknowledged = make(map[nodes.Address]time.Time, len(rec.Acknowledged))
	for k, v := range rec.Acknowledged {
		cp.Acknowledged[k] = v
	}
	return cp, true
}

// Pending returns replication records that are not yet completed.
func (r *Replicator) Pending() []ReplicationRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]ReplicationRecord, 0)
	for _, rec := range r.records {
		if rec.Completed {
			continue
		}
		cp := ReplicationRecord{
			Hash:        rec.Hash,
			Attempts:    rec.Attempts,
			LastAttempt: rec.LastAttempt,
			Completed:   rec.Completed,
		}
		cp.Acknowledged = make(map[nodes.Address]time.Time, len(rec.Acknowledged))
		for k, v := range rec.Acknowledged {
			cp.Acknowledged[k] = v
		}
		out = append(out, cp)
	}
	return out
}

// SetRetryPolicy adjusts retry attempts and delay.
func (r *Replicator) SetRetryPolicy(limit int, delay time.Duration) {
	r.mu.Lock()
	if limit >= 0 {
		r.retryLimit = limit
	}
	if delay > 0 {
		r.retryDelay = delay
	}
	r.mu.Unlock()
}

// SetEnqueueTimeout sets how long replication requests wait for queue space.
func (r *Replicator) SetEnqueueTimeout(timeout time.Duration) {
	r.mu.Lock()
	r.timeout = timeout
	r.mu.Unlock()
}

// Metrics returns counters describing replication progress.
func (r *Replicator) Metrics() ReplicationMetrics {
	pending := r.Pending()
	return ReplicationMetrics{
		Enqueued:     r.metrics.enqueued.Load(),
		Acknowledged: r.metrics.acknowledged.Load(),
		Retries:      r.metrics.retries.Load(),
		Dropped:      r.metrics.dropped.Load(),
		Pending:      len(pending),
	}
}
