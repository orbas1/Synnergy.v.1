package core

import "testing"

func TestRollupManager(t *testing.T) {
	agg := NewRollupAggregator()
	mgr := NewRollupManager(agg)
	mgr.Pause()
	if !mgr.Status() {
		t.Fatalf("expected paused")
	}
	mgr.Resume()
	if mgr.Status() {
		t.Fatalf("expected running")
	}

	if _, err := mgr.Submit([]string{"tx1"}); err != nil {
		t.Fatalf("submit: %v", err)
	}

	empty := NewRollupManager(nil)
	if _, err := empty.Submit([]string{"tx1"}); err == nil {
		t.Fatalf("expected error for nil aggregator")
	}
}
