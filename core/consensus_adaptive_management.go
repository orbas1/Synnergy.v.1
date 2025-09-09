package core

import "sync"

// AdaptiveManager tunes consensus weights based on network metrics.
type AdaptiveManager struct {
	mu      sync.RWMutex
	engine  *SynnergyConsensus
	window  int
	metrics []netMetric
}

type netMetric struct {
	demand float64
	stake  float64
}

// NewAdaptiveManager creates a manager for the given consensus engine. An optional
// window size controls how many recent metrics are averaged; defaults to 10.
func NewAdaptiveManager(engine *SynnergyConsensus, window ...int) *AdaptiveManager {
	w := 10
	if len(window) > 0 && window[0] > 0 {
		w = window[0]
	}
	return &AdaptiveManager{engine: engine, window: w}
}

// record appends a metric sample and maintains the sliding window.
func (am *AdaptiveManager) record(demand, stake float64) {
	am.metrics = append(am.metrics, netMetric{demand: demand, stake: stake})
	if len(am.metrics) > am.window {
		am.metrics = am.metrics[1:]
	}
}

func (am *AdaptiveManager) averages() (d, s float64) {
	for _, m := range am.metrics {
		d += m.demand
		s += m.stake
	}
	n := float64(len(am.metrics))
	if n == 0 {
		return 0, 0
	}
	return d / n, s / n
}

// Adjust recalculates weights using the provided demand and stake
// concentrations. It returns the updated weights for convenience.
func (am *AdaptiveManager) Adjust(demand, stake float64) ConsensusWeights {
	am.mu.Lock()
	defer am.mu.Unlock()
	if am.engine == nil {
		return ConsensusWeights{}
	}
	am.record(demand, stake)
	d, s := am.averages()
	am.engine.AdjustWeights(d, s)
	return am.engine.Weights
}

// Threshold computes the consensus threshold for switching modes using the
// supplied metrics.
func (am *AdaptiveManager) Threshold(demand, stake float64) float64 {
	am.mu.Lock()
	defer am.mu.Unlock()
	if am.engine == nil {
		return 0
	}
	am.record(demand, stake)
	d, s := am.averages()
	return am.engine.Threshold(d, s)
}

// Weights returns the current weights without modifying them.
func (am *AdaptiveManager) Weights() ConsensusWeights {
	am.mu.RLock()
	defer am.mu.RUnlock()
	if am.engine == nil {
		return ConsensusWeights{}
	}
	return am.engine.Weights
}

// RecordMetrics allows external components to feed network demand and stake
// observations without immediately adjusting weights.
func (am *AdaptiveManager) RecordMetrics(demand, stake float64) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.record(demand, stake)
}

// Reset clears all recorded metrics.
func (am *AdaptiveManager) Reset() {
	am.mu.Lock()
	am.metrics = nil
	am.mu.Unlock()
}
