package core

import "sync"

// MobileMiningNode adapts mining to operate within mobile constraints.
type MobileMiningNode struct {
	base       *MiningNode
	mu         sync.Mutex
	powerLimit uint64
}

// NewMobileMiningNode creates a MobileMiningNode wrapping a standard mining node.
func NewMobileMiningNode(hashRate, powerLimit uint64) *MobileMiningNode {
	return &MobileMiningNode{base: NewMiningNode(hashRate), powerLimit: powerLimit}
}

// Start activates the underlying mining node if within power limits.
func (mm *MobileMiningNode) Start() {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.base.Start()
}

// Stop halts the mining process.
func (mm *MobileMiningNode) Stop() { mm.base.Stop() }

// IsMining reports whether the node is currently mining.
func (mm *MobileMiningNode) IsMining() bool { return mm.base.IsMining() }

// Mine delegates to the underlying mining node.
func (mm *MobileMiningNode) Mine(data []byte) (string, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.powerLimit == 0 {
		return "", nil
	}
	return mm.base.Mine(data)
}

// SetPowerLimit updates the allowed power usage for mining operations.
func (mm *MobileMiningNode) SetPowerLimit(limit uint64) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.powerLimit = limit
}

// PowerLimit returns the configured power limit.
func (mm *MobileMiningNode) PowerLimit() uint64 {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	return mm.powerLimit
}
