package synnergy

import (
	"math"
	"sync"
)

// AnomalyDetector performs streaming anomaly detection using mean and variance.
type AnomalyDetector struct {
	mu        sync.RWMutex
	count     float64
	mean      float64
	m2        float64
	threshold float64
}

// NewAnomalyDetector constructs a detector with a z-score threshold.
// If a non-positive threshold is supplied, a default of 3 is used to avoid
// disabling anomaly detection inadvertently.
func NewAnomalyDetector(threshold float64) *AnomalyDetector {
	if threshold <= 0 {
		threshold = 3
	}
	return &AnomalyDetector{threshold: threshold}
}

// Update incorporates a new observation into the running statistics.
func (a *AnomalyDetector) Update(v float64) {
	a.mu.Lock()
	a.count++
	delta := v - a.mean
	a.mean += delta / a.count
	a.m2 += delta * (v - a.mean)
	a.mu.Unlock()
}

// IsAnomalous reports whether the value deviates beyond the configured threshold.
func (a *AnomalyDetector) IsAnomalous(v float64) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.count < 2 {
		return false
	}
	variance := a.m2 / (a.count - 1)
	if variance == 0 {
		return false
	}
	z := math.Abs((v - a.mean) / math.Sqrt(variance))
	return z > a.threshold
}
