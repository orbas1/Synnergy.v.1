package nodes

import "time"

// BlockSummary provides the minimal metadata required to reference a
// block in archival storage.
type BlockSummary struct {
	Height    uint64
	Hash      string
	Timestamp time.Time
}

// HistoricalNodeInterface extends NodeInterface with archival
// functionality for serving historical chain data.
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
