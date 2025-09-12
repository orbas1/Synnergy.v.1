package core

import (
	"testing"
)

func TestSystemHealthLogger(t *testing.T) {
	l := NewSystemHealthLogger()
	m := l.Collect(-1, 10)
	if m.PeerCount != 0 {
		t.Fatalf("expected peer count clamped to 0 got %d", m.PeerCount)
	}
	snap := l.Snapshot()
	if snap.LastBlockHeight != 10 {
		t.Fatalf("snapshot mismatch")
	}
}
