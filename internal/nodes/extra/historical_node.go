package nodes

import (
	"errors"
	"time"
)

// BlockSummary provides the minimal metadata required to reference a block in
// archival storage.
type BlockSummary struct {
	Height    uint64
	Hash      string
	Timestamp time.Time
}

// ArchivePolicy configures how historical nodes retain data.
type ArchivePolicy struct {
	MaxEntries int
	Retention  time.Duration
}

// HistoricalOption configures optional behaviour for an archival node.
type HistoricalOption interface {
	applyHistoricalOption(*ArchivePolicy)
}

type historicalOptionFunc func(*ArchivePolicy)

func (f historicalOptionFunc) applyHistoricalOption(p *ArchivePolicy) { f(p) }

// WithArchiveLimit constrains the maximum number of entries retained.
func WithArchiveLimit(limit int) HistoricalOption {
	return historicalOptionFunc(func(p *ArchivePolicy) {
		if limit > 0 {
			p.MaxEntries = limit
		}
	})
}

// WithArchiveRetention sets the retention window for archived blocks.
func WithArchiveRetention(window time.Duration) HistoricalOption {
	return historicalOptionFunc(func(p *ArchivePolicy) {
		if window > 0 {
			p.Retention = window
		}
	})
}

// ApplyHistoricalOptions applies the supplied options to a base policy.
func ApplyHistoricalOptions(base ArchivePolicy, opts ...HistoricalOption) ArchivePolicy {
	policy := base
	for _, opt := range opts {
		if opt != nil {
			opt.applyHistoricalOption(&policy)
		}
	}
	return policy
}

// ValidateBlockSummary ensures a block summary contains the minimal required
// data before being stored.
func ValidateBlockSummary(summary BlockSummary) error {
	if summary.Hash == "" {
		return errors.New("block hash cannot be empty")
	}
	if summary.Timestamp.IsZero() {
		return errors.New("timestamp must be set")
	}
	return nil
}

// HistoricalNodeInterface extends NodeInterface with archival functionality for
// serving historical chain data.
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
