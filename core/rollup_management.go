package core

import "errors"

// RollupManager provides administrative controls for a RollupAggregator.
type RollupManager struct {
	agg *RollupAggregator
}

// NewRollupManager binds management utilities to an aggregator.
func NewRollupManager(agg *RollupAggregator) *RollupManager {
	return &RollupManager{agg: agg}
}

// Pause halts the underlying aggregator.
func (m *RollupManager) Pause() {
	if m.agg != nil {
		m.agg.Pause()
	}
}

// Resume restarts the aggregator.
func (m *RollupManager) Resume() {
	if m.agg != nil {
		m.agg.Resume()
	}
}

// Status returns whether the aggregator is paused.
func (m *RollupManager) Status() bool {
	if m.agg == nil {
		return false
	}
	return m.agg.Status()
}

// Submit delegates to the underlying aggregator if available.
func (m *RollupManager) Submit(txs []string) (string, error) {
	if m.agg == nil {
		return "", errors.New("no aggregator")
	}
	return m.agg.SubmitBatch(txs)
}
