package core

import (
	"sync"
)

// ShardManager maintains shard leaders and cross-shard transaction receipts.
type ShardManager struct {
	mu         sync.RWMutex
	leaders    map[int]string
	receipts   map[int][]string // shardID -> tx hashes destined for shard
	shardBits  uint8
	shardLoads map[int]int
}

// NewShardManager creates an empty shard manager.
func NewShardManager(shardBits uint8) *ShardManager {
	return &ShardManager{
		leaders:    make(map[int]string),
		receipts:   make(map[int][]string),
		shardBits:  shardBits,
		shardLoads: make(map[int]int),
	}
}

// GetLeader returns the leader address for a shard.
func (m *ShardManager) GetLeader(shardID int) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	addr, ok := m.leaders[shardID]
	return addr, ok
}

// SetLeader sets the leader address for a shard.
func (m *ShardManager) SetLeader(shardID int, addr string) {
	m.mu.Lock()
	m.leaders[shardID] = addr
	m.mu.Unlock()
}

// LeaderMap returns the mapping of shards to leaders.
func (m *ShardManager) LeaderMap() map[int]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[int]string, len(m.leaders))
	for id, addr := range m.leaders {
		out[id] = addr
	}
	return out
}

// SubmitCrossShardTx records a cross-shard transaction header.
func (m *ShardManager) SubmitCrossShardTx(fromShard, toShard int, txHash string) {
	m.mu.Lock()
	m.receipts[toShard] = append(m.receipts[toShard], txHash)
	m.shardLoads[toShard]++
	m.mu.Unlock()
}

// PullReceipts retrieves receipts for a shard and clears them.
func (m *ShardManager) PullReceipts(shardID int) []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	txs := m.receipts[shardID]
	m.receipts[shardID] = nil
	m.shardLoads[shardID] = 0
	return txs
}

// Reshard updates the shard bit-size (i.e., number of shards = 2^bits).
func (m *ShardManager) Reshard(newBits uint8) {
	m.mu.Lock()
	m.shardBits = newBits
	m.mu.Unlock()
}

// Rebalance returns shard IDs exceeding the load threshold.
func (m *ShardManager) Rebalance(threshold int) []int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var heavy []int
	for id, load := range m.shardLoads {
		if load > threshold {
			heavy = append(heavy, id)
		}
	}
	return heavy
}

// ShardCount returns the current number of shards.
func (m *ShardManager) ShardCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return 1 << m.shardBits
}
