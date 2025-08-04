package synnergy

import (
	"errors"
	"sync"
)

// MobileMiningNode wraps a MiningNode with additional resource management suited
// for battery powered devices.  Mining automatically pauses when the battery
// level drops below a configurable threshold.
type MobileMiningNode struct {
	*MiningNode
	mu        sync.RWMutex
	battery   float64
	threshold float64
}

// NewMobileMiningNode constructs a mobile miner with an initial battery level
// and threshold at which mining should pause.
func NewMobileMiningNode(id string, hashRate, battery, threshold float64) *MobileMiningNode {
	return &MobileMiningNode{MiningNode: NewMiningNode(id, hashRate), battery: battery, threshold: threshold}
}

// UpdateBattery records the latest battery level for the device.
func (m *MobileMiningNode) UpdateBattery(level float64) {
	m.mu.Lock()
	m.battery = level
	if m.battery < m.threshold {
		m.MiningNode.Stop()
	}
	m.mu.Unlock()
}

// Start begins mining if the battery level is above the configured threshold.
func (m *MobileMiningNode) Start() error {
	m.mu.RLock()
	ok := m.battery >= m.threshold
	m.mu.RUnlock()
	if !ok {
		return errors.New("battery level too low for mining")
	}
	m.MiningNode.Start()
	return nil
}

// Battery returns the current battery level.
func (m *MobileMiningNode) Battery() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.battery
}

// Threshold returns the minimum battery level required to mine.
func (m *MobileMiningNode) Threshold() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.threshold
}

// SetThreshold updates the minimum battery level required to mine.
func (m *MobileMiningNode) SetThreshold(th float64) {
	m.mu.Lock()
	m.threshold = th
	m.mu.Unlock()
}
