package synnergy

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// MiningNode represents a basic proof-of-work miner that can generate block
// candidates and submit mined blocks to the network.  The node intentionally
// keeps state lightweight so that it can operate independently from the rest of
// the system for tests or simulations.
type MiningNode struct {
	id        string
	mu        sync.RWMutex
	running   bool
	hashRate  float64 // hashes per second
	lastBlock string
}

// NewMiningNode constructs a mining node with the supplied identifier and
// nominal hash rate.
func NewMiningNode(id string, hashRate float64) *MiningNode {
	return &MiningNode{id: id, hashRate: hashRate}
}

// ID returns the miner identifier.
func (m *MiningNode) ID() string {
	return m.id
}

// Start begins the mining process.  It is safe to call Start multiple times; the
// mining loop will only run once until Stop is invoked.
func (m *MiningNode) Start() {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return
	}
	m.running = true
	m.mu.Unlock()

	go m.mineLoop()
}

// Stop signals the miner to halt work.
func (m *MiningNode) Stop() {
	m.mu.Lock()
	m.running = false
	m.mu.Unlock()
}

// IsRunning reports whether the miner is actively hashing.
func (m *MiningNode) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// HashRate returns the configured hash rate for the miner.
func (m *MiningNode) HashRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.hashRate
}

// SubmitBlock records the hash of a mined block.  In a full implementation this
// would broadcast the block to peers.  Here we simply store the hash for tests.
func (m *MiningNode) SubmitBlock(hash string) {
	m.mu.Lock()
	m.lastBlock = hash
	m.mu.Unlock()
}

// LastBlock returns the most recently submitted block hash.
func (m *MiningNode) LastBlock() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastBlock
}

// mineLoop simulates mining by periodically generating random block hashes at
// the configured hash rate.
func (m *MiningNode) mineLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		m.mu.RLock()
		running := m.running
		rate := m.hashRate
		m.mu.RUnlock()
		if !running {
			return
		}
		hashes := int(rate)
		for i := 0; i < hashes; i++ {
			// simulate a mined block
			hash := randomHash()
			m.SubmitBlock(hash)
		}
	}
}

// randomHash is a helper that returns a pseudo-random hexadecimal string.
func randomHash() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// MineBlock performs a synchronous mining attempt with the given difficulty.
// A valid block hash with a number of leading zero bits equal to difficulty will
// be returned.  This method is primarily intended for tests and does not perform
// any network operations.
func (m *MiningNode) MineBlock(difficulty int) (string, error) {
	if difficulty < 0 || difficulty > 256 {
		return "", errors.New("invalid difficulty")
	}
	prefix := strings.Repeat("0", difficulty/4)
	for {
		hash := randomHash()
		if strings.HasPrefix(hash, prefix) {
			m.SubmitBlock(hash)
			return hash, nil
		}
	}
}
