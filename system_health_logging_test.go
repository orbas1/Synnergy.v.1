package synnergy

import "testing"

func TestSystemHealthLogger(t *testing.T) {
	l := NewSystemHealthLogger()
	m := l.Collect(5, 10)
	if m.PeerCount != 5 || m.LastBlockHeight != 10 {
		t.Fatalf("unexpected metrics: %+v", m)
	}
	snap := l.Snapshot()
	if snap.PeerCount != 5 || snap.LastBlockHeight != 10 {
		t.Fatalf("snapshot mismatch: %+v", snap)
	}
}
