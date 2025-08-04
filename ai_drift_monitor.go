package synnergy

import (
	"math"
	"sync"
)

// DriftMonitor tracks model performance metrics to detect drift.
type DriftMonitor struct {
	mu        sync.RWMutex
	baselines map[string]float64
}

// NewDriftMonitor creates a new drift monitor.
func NewDriftMonitor() *DriftMonitor {
	return &DriftMonitor{baselines: make(map[string]float64)}
}

// UpdateBaseline sets the baseline metric for a model.
func (d *DriftMonitor) UpdateBaseline(modelHash string, metric float64) {
	d.mu.Lock()
	d.baselines[modelHash] = metric
	d.mu.Unlock()
}

// HasDrift reports if the new metric deviates from baseline by more than threshold.
func (d *DriftMonitor) HasDrift(modelHash string, metric, threshold float64) bool {
	d.mu.RLock()
	base, ok := d.baselines[modelHash]
	d.mu.RUnlock()
	if !ok {
		return false
	}
	return math.Abs(metric-base) > threshold
}
