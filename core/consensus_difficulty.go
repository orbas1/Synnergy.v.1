package core

import "sync"

// DifficultyManager maintains PoW difficulty using a sliding window of
// recent block times.
type DifficultyManager struct {
	mu         sync.Mutex
	engine     *SynnergyConsensus
	window     int
	target     float64
	difficulty float64
	samples    []float64
}

// NewDifficultyManager initialises the manager.
func NewDifficultyManager(engine *SynnergyConsensus, window int, initial, target float64) *DifficultyManager {
	if window <= 0 {
		window = 1
	}
	return &DifficultyManager{engine: engine, window: window, target: target, difficulty: initial}
}

// AddSample records the duration of the last block and recomputes difficulty.
// It returns the updated difficulty value.
func (dm *DifficultyManager) AddSample(duration float64) float64 {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.samples = append(dm.samples, duration)
	if len(dm.samples) > dm.window {
		dm.samples = dm.samples[1:]
	}
	var total float64
	for _, d := range dm.samples {
		total += d
	}
	avg := total / float64(len(dm.samples))
	if dm.engine != nil {
		dm.difficulty = dm.engine.DifficultyAdjust(dm.difficulty, avg, dm.target)
	}
	return dm.difficulty
}

// Difficulty returns the current difficulty value.
func (dm *DifficultyManager) Difficulty() float64 {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	return dm.difficulty
}
