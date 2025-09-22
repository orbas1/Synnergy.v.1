package nodes

import (
	"testing"
	"time"
)

func TestApplyHistoricalOptions(t *testing.T) {
	policy := ApplyHistoricalOptions(ArchivePolicy{}, WithArchiveLimit(10), WithArchiveRetention(time.Hour))
	if policy.MaxEntries != 10 {
		t.Fatalf("expected max entries 10 got %d", policy.MaxEntries)
	}
	if policy.Retention != time.Hour {
		t.Fatalf("expected retention hour got %v", policy.Retention)
	}
}

func TestValidateBlockSummary(t *testing.T) {
	if err := ValidateBlockSummary(BlockSummary{}); err == nil {
		t.Fatalf("expected error for empty summary")
	}
	summary := BlockSummary{Height: 1, Hash: "hash", Timestamp: time.Now()}
	if err := ValidateBlockSummary(summary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
