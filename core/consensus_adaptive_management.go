package core

import "sync"

// AdaptiveManager tunes consensus weights based on network metrics.
type AdaptiveManager struct {
	mu     sync.RWMutex
	engine *SynnergyConsensus
}

// NewAdaptiveManager creates a manager for the given consensus engine.
func NewAdaptiveManager(engine *SynnergyConsensus) *AdaptiveManager {
	return &AdaptiveManager{engine: engine}
}

// Adjust recalculates weights using the provided demand and stake
// concentrations. It returns the updated weights for convenience.
func (am *AdaptiveManager) Adjust(demand, stake float64) ConsensusWeights {
	am.mu.Lock()
	defer am.mu.Unlock()
	if am.engine == nil {
		return ConsensusWeights{}
	}
	am.engine.AdjustWeights(demand, stake)
	return am.engine.Weights
}

// Threshold computes the consensus threshold for switching modes using the
// supplied metrics.
func (am *AdaptiveManager) Threshold(demand, stake float64) float64 {
	am.mu.RLock()
	defer am.mu.RUnlock()
	if am.engine == nil {
		return 0
	}
	return am.engine.Threshold(demand, stake)
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
