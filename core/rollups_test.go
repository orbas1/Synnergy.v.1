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
}
