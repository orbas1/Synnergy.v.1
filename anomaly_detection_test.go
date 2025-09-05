package synnergy

import "testing"

func TestAnomalyDetectorBasic(t *testing.T) {
	d := NewAnomalyDetector(2)
	values := []float64{10, 11, 9, 10, 11, 9}
	for _, v := range values {
		d.Update(v)
	}
	if d.IsAnomalous(10) {
		t.Fatalf("10 should not be anomalous")
	}
	if !d.IsAnomalous(30) {
		t.Fatalf("30 should be anomalous")
	}
}

func TestAnomalyDetectorDefaultThreshold(t *testing.T) {
	d := NewAnomalyDetector(0)
	if d.threshold != 3 {
		t.Fatalf("expected default threshold 3, got %f", d.threshold)
	}
}
