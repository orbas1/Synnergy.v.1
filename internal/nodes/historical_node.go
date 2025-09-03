package nodes

import (
	"sync"
	"time"
)

// BlockSummary provides minimal metadata required to reference a block in
// archival storage.
type BlockSummary struct {
	Height    uint64
	Hash      string
	Timestamp time.Time
}

// HistoricalNodeInterface defines archival functionality for serving historical
// chain data.
type HistoricalNodeInterface interface {
	// ArchiveBlock stores a block summary for long term retrieval.
	ArchiveBlock(summary BlockSummary) error
	// GetBlockByHeight retrieves a stored block by its height.
	GetBlockByHeight(height uint64) (BlockSummary, bool)
	// GetBlockByHash retrieves a stored block by its hash.
	GetBlockByHash(hash string) (BlockSummary, bool)
	// TotalBlocks returns the number of archived blocks.
	TotalBlocks() int
}

// HistoricalNode provides an in-memory archive of block summaries backed by a
// BasicNode. The implementation is intended for simulations and does not
// persist data to disk.
type HistoricalNode struct {
	*BasicNode
	mu       sync.RWMutex
	byHeight map[uint64]BlockSummary
	byHash   map[string]BlockSummary
}

// NewHistoricalNode creates a new historical node.
func NewHistoricalNode(id Address) *HistoricalNode {
	return &HistoricalNode{BasicNode: NewBasicNode(id), byHeight: make(map[uint64]BlockSummary), byHash: make(map[string]BlockSummary)}
}

// ArchiveBlock stores a block summary for later retrieval.
func (n *HistoricalNode) ArchiveBlock(summary BlockSummary) error {
	n.mu.Lock()
	n.byHeight[summary.Height] = summary
	n.byHash[summary.Hash] = summary
	n.mu.Unlock()
	return nil
}

// GetBlockByHeight retrieves a stored block by its height.
func (n *HistoricalNode) GetBlockByHeight(height uint64) (BlockSummary, bool) {
	n.mu.RLock()
	bs, ok := n.byHeight[height]
	n.mu.RUnlock()
	return bs, ok
}

// GetBlockByHash retrieves a stored block by its hash.
func (n *HistoricalNode) GetBlockByHash(hash string) (BlockSummary, bool) {
	n.mu.RLock()
	bs, ok := n.byHash[hash]
	n.mu.RUnlock()
	return bs, ok
}

// TotalBlocks returns the number of archived blocks.
func (n *HistoricalNode) TotalBlocks() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return len(n.byHeight)
}

// Ensure HistoricalNode implements HistoricalNodeInterface.
var _ HistoricalNodeInterface = (*HistoricalNode)(nil)
