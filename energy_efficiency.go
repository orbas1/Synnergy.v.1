package synnergy

import "sync"

// EfficiencyRecord stores energy usage and processed transactions for a validator.
type EfficiencyRecord struct {
	Transactions int
	EnergyKWh    float64
}

// EnergyEfficiencyTracker records validator efficiency metrics.
type EnergyEfficiencyTracker struct {
	mu    sync.RWMutex
	stats map[string]EfficiencyRecord
}

// NewEnergyEfficiencyTracker constructs a tracker instance.
func NewEnergyEfficiencyTracker() *EnergyEfficiencyTracker {
	return &EnergyEfficiencyTracker{stats: make(map[string]EfficiencyRecord)}
}

// Record stores processed transaction and energy usage values for a validator.
func (t *EnergyEfficiencyTracker) Record(validator string, txProcessed int, energyKWh float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	r := t.stats[validator]
	r.Transactions += txProcessed
	r.EnergyKWh += energyKWh
	t.stats[validator] = r
}

// Efficiency returns transactions per kilowatt hour for a validator.
// The second return value reports whether metrics were found.
func (t *EnergyEfficiencyTracker) Efficiency(validator string) (float64, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	r, ok := t.stats[validator]
	if !ok || r.EnergyKWh == 0 {
		return 0, ok
	}
	return float64(r.Transactions) / r.EnergyKWh, true
}

// NetworkAverage calculates the average transactions per kWh across all validators.
func (t *EnergyEfficiencyTracker) NetworkAverage() float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var totalTx int
	var totalEnergy float64
	for _, r := range t.stats {
		totalTx += r.Transactions
		totalEnergy += r.EnergyKWh
	}
	if totalEnergy == 0 {
		return 0
	}
	return float64(totalTx) / totalEnergy
}

// Stats returns the raw efficiency record for a validator.
// The boolean result reports whether metrics were recorded for the validator.
func (t *EnergyEfficiencyTracker) Stats(validator string) (EfficiencyRecord, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	r, ok := t.stats[validator]
	return r, ok
}

// Reset removes all tracked metrics for the given validator.
func (t *EnergyEfficiencyTracker) Reset(validator string) {
	t.mu.Lock()
	delete(t.stats, validator)
	t.mu.Unlock()
}
