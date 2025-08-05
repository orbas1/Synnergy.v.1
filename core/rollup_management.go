package core

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
	m.agg.Pause()
}

// Resume restarts the aggregator.
func (m *RollupManager) Resume() {
	m.agg.Resume()
}

// Status returns whether the aggregator is paused.
func (m *RollupManager) Status() bool {
	return m.agg.Status()
}
