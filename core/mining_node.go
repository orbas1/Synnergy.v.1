package core

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
)

// MiningNode performs proof-of-work style mining operations.
type MiningNode struct {
	mu       sync.Mutex
	active   bool
	hashRate uint64
	nonce    uint64
}

// NewMiningNode creates a mining node with the supplied hashrate hint.
func NewMiningNode(hashRate uint64) *MiningNode {
	return &MiningNode{hashRate: hashRate}
}

// Start activates the mining process.
func (mn *MiningNode) Start() {
	mn.mu.Lock()
	mn.active = true
	mn.mu.Unlock()
}

// Stop deactivates mining operations.
func (mn *MiningNode) Stop() {
	mn.mu.Lock()
	mn.active = false
	mn.mu.Unlock()
}

// IsMining reports whether the node is currently mining.
func (mn *MiningNode) IsMining() bool {
	mn.mu.Lock()
	defer mn.mu.Unlock()
	return mn.active
}

// Mine runs a single hashing attempt over data and returns the resulting hash.
func (mn *MiningNode) Mine(data []byte) (string, error) {
	mn.mu.Lock()
	if !mn.active {
		mn.mu.Unlock()
		return "", errors.New("mining inactive")
	}
	nonce := mn.nonce
	mn.nonce++
	mn.mu.Unlock()
	buf := append(data, byte(nonce))
	h := sha256.Sum256(buf)
	return hex.EncodeToString(h[:]), nil
}

// HashRateHint returns the configured hash rate for the node.
func (mn *MiningNode) HashRateHint() uint64 {
	mn.mu.Lock()
	defer mn.mu.Unlock()
	return mn.hashRate
}
