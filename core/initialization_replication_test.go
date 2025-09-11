package core

import "testing"

func TestInitServiceStartStop(t *testing.T) {
	l := NewLedger()
	r := NewReplicator(l)
	init := NewInitService(r)
	if err := init.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if !r.Status() || !init.Status() {
		t.Fatalf("service should be running")
	}
	if err := init.Start(); err != ErrInitServiceRunning {
		t.Fatalf("expected running error")
	}
	if err := init.Stop(); err != nil {
		t.Fatalf("stop failed: %v", err)
	}
	if r.Status() || init.Status() {
		t.Fatalf("service should be stopped")
	}
	if err := init.Stop(); err != ErrInitServiceStopped {
		t.Fatalf("expected stopped error")
	}
}

func TestInitServiceNilReplicator(t *testing.T) {
	init := NewInitService(nil)
	if err := init.Start(); err != ErrInitReplicatorNil {
		t.Fatalf("expected nil replicator error")
	}
}
