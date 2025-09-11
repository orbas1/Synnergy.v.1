package core

import "testing"

func TestRollupAggregator(t *testing.T) {
	agg := NewRollupAggregator()
	id, err := agg.SubmitBatch([]string{"tx1", "tx2"})
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if id == "" {
		t.Fatalf("expected batch id")
	}
	txs, err := agg.BatchTransactions(id)
	if err != nil || len(txs) != 2 {
		t.Fatalf("unexpected transactions")
	}
	agg.Pause()
	if !agg.Status() {
		t.Fatalf("expected paused")
	}
	agg.Resume()
	if agg.Status() {
		t.Fatalf("expected resumed")
	}

	if err := agg.FinalizeBatch(id, true); err != nil {
		t.Fatalf("finalize: %v", err)
	}
	finalized := agg.ListBatchesByStatus("finalized")
	if len(finalized) != 1 || finalized[0].ID != id {
		t.Fatalf("expected finalized batch in list")
	}
}
