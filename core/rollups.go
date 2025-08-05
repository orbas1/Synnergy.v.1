package core

import (
	"errors"
	"fmt"
	"sync"
)

// RollupBatch represents a batch of L2 transactions submitted to L1.
type RollupBatch struct {
	ID           string
	Transactions []string
	Status       string // pending, challenged, finalized, reverted
}

// RollupAggregator coordinates rollup batch submission and state.
type RollupAggregator struct {
	mu      sync.RWMutex
	seq     int
	batches map[string]*RollupBatch
	paused  bool
}

// NewRollupAggregator constructs an empty RollupAggregator.
func NewRollupAggregator() *RollupAggregator {
	return &RollupAggregator{batches: make(map[string]*RollupBatch)}
}

// SubmitBatch records a new rollup batch. Returns the batch ID.
func (a *RollupAggregator) SubmitBatch(txs []string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.paused {
		return "", errors.New("rollup aggregator paused")
	}
	a.seq++
	id := fmt.Sprintf("batch-%d", a.seq)
	a.batches[id] = &RollupBatch{ID: id, Transactions: append([]string(nil), txs...), Status: "pending"}
	return id, nil
}

// ChallengeBatch marks a batch as challenged.
func (a *RollupAggregator) ChallengeBatch(id string, txIdx int, proof []byte) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	b, ok := a.batches[id]
	if !ok {
		return fmt.Errorf("batch %s not found", id)
	}
	if txIdx < 0 || txIdx >= len(b.Transactions) {
		return fmt.Errorf("invalid transaction index")
	}
	// Proof verification is out of scope for this lightweight module.
	b.Status = "challenged"
	return nil
}

// FinalizeBatch finalizes a batch as valid or reverts it.
func (a *RollupAggregator) FinalizeBatch(id string, valid bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	b, ok := a.batches[id]
	if !ok {
		return fmt.Errorf("batch %s not found", id)
	}
	if valid {
		b.Status = "finalized"
	} else {
		b.Status = "reverted"
	}
	return nil
}

// BatchInfo retrieves a batch by ID.
func (a *RollupAggregator) BatchInfo(id string) (RollupBatch, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	b, ok := a.batches[id]
	if !ok {
		return RollupBatch{}, false
	}
	return *b, true
}

// ListBatches lists all recorded batches.
func (a *RollupAggregator) ListBatches() []RollupBatch {
	a.mu.RLock()
	defer a.mu.RUnlock()
	res := make([]RollupBatch, 0, len(a.batches))
	for _, b := range a.batches {
		res = append(res, *b)
	}
	return res
}

// BatchTransactions returns transactions included in a batch.
func (a *RollupAggregator) BatchTransactions(id string) ([]string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	b, ok := a.batches[id]
	if !ok {
		return nil, fmt.Errorf("batch %s not found", id)
	}
	return append([]string(nil), b.Transactions...), nil
}

// Pause halts new batch submissions.
func (a *RollupAggregator) Pause() {
	a.mu.Lock()
	a.paused = true
	a.mu.Unlock()
}

// Resume allows batch submissions.
func (a *RollupAggregator) Resume() {
	a.mu.Lock()
	a.paused = false
	a.mu.Unlock()
}

// Status reports whether the aggregator is paused.
func (a *RollupAggregator) Status() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.paused
}
