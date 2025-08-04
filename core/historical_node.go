package core

import (
	"errors"
	"sync"
	nodes "synnergy/Nodes"
)

// HistoricalNode provides archival functionality for serving historical
// chain data. It keeps block summaries in memory for quick lookup.
type HistoricalNode struct {
	mu       sync.RWMutex
	byHeight map[uint64]nodes.BlockSummary
	byHash   map[string]nodes.BlockSummary
}

// Ensure HistoricalNode implements nodes.HistoricalNodeInterface.
var _ nodes.HistoricalNodeInterface = (*HistoricalNode)(nil)

// NewHistoricalNode creates an empty HistoricalNode.
func NewHistoricalNode() *HistoricalNode {
	return &HistoricalNode{
		byHeight: make(map[uint64]nodes.BlockSummary),
		byHash:   make(map[string]nodes.BlockSummary),
	}
}

// ArchiveBlock stores a block summary for long term retrieval.
func (h *HistoricalNode) ArchiveBlock(summary nodes.BlockSummary) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.byHeight[summary.Height]; exists {
		return errors.New("block height already archived")
	}
	if _, exists := h.byHash[summary.Hash]; exists {
		return errors.New("block hash already archived")
	}
	h.byHeight[summary.Height] = summary
	h.byHash[summary.Hash] = summary
	return nil
}

// GetBlockByHeight retrieves a stored block by its height.
func (h *HistoricalNode) GetBlockByHeight(height uint64) (nodes.BlockSummary, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	summary, ok := h.byHeight[height]
	return summary, ok
}

// GetBlockByHash retrieves a stored block by its hash.
func (h *HistoricalNode) GetBlockByHash(hash string) (nodes.BlockSummary, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	summary, ok := h.byHash[hash]
	return summary, ok
}

// TotalBlocks returns the number of archived blocks.
func (h *HistoricalNode) TotalBlocks() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.byHeight)
}
