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
}
