package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
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

// MineUntil keeps hashing the provided data until the resulting hash has the
// given prefix or the context is cancelled. The resulting hash and the nonce
// used to produce it are returned. A non-empty prefix must be supplied.
func (mn *MiningNode) MineUntil(ctx context.Context, data []byte, prefix string) (string, uint64, error) {
	if prefix == "" {
		return "", 0, errors.New("prefix required")
	}
	for {
		select {
		case <-ctx.Done():
			return "", 0, ctx.Err()
		default:
			hash, err := mn.Mine(data)
			if err != nil {
				return "", 0, err
			}
			if strings.HasPrefix(hash, prefix) {
				mn.mu.Lock()
				nonce := mn.nonce
				mn.mu.Unlock()
				return hash, nonce, nil
			}
		}
	}
}
