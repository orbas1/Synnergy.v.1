package core

import "testing"

func TestSyncManagerLifecycle(t *testing.T) {
	l := NewLedger()
	sm := NewSyncManager(l)

	// should error when not running
	if err := sm.Once(); err == nil {
		t.Fatalf("expected error when not running")
	}

	sm.Start()
	running, _ := sm.Status()
	if !running {
		t.Fatalf("expected running")
	}

	l.AddBlock(&Block{Hash: "b1"})
	if err := sm.Once(); err != nil {
		t.Fatalf("once: %v", err)
	}
	_, h := sm.Status()
	if h != 1 {
		t.Fatalf("expected height 1 got %d", h)
	}

	sm.Stop()
	running, _ = sm.Status()
	if running {
		t.Fatalf("expected stopped")
	}
	if err := sm.Once(); err == nil {
		t.Fatalf("expected error when stopped")
	}
}
