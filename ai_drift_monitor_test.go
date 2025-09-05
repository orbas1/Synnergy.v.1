package synnergy

import "testing"

func TestDriftMonitorBaseline(t *testing.T) {
	d := NewDriftMonitor()
	d.UpdateBaseline("m1", 0.5)
	if !d.HasDrift("m1", 0.8, 0.2) {
		t.Fatalf("expected drift to be detected")
	}
	if d.HasDrift("m1", 0.55, 0.2) {
		t.Fatalf("unexpected drift detected")
	}
	if d.HasDrift("unknown", 0.8, 0.2) {
		t.Fatalf("unknown model should not report drift")
	}
}
