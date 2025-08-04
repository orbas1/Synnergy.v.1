package core

import "testing"

func TestInitServiceStartStop(t *testing.T) {
	l := NewLedger()
	r := NewReplicator(l)
	init := NewInitService(r)

	init.Start()
	if !r.Status() {
		t.Fatalf("replicator should be running")
	}

	init.Stop()
	if r.Status() {
		t.Fatalf("replicator should be stopped")
	}
}
