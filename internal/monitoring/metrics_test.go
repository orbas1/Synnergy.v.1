package monitoring

import "testing"

func TestMetrics(t *testing.T) {
	m := NewMetrics()
	m.Inc("a")
	if m.Get("a") != 1 {
		t.Fatalf("expected 1")
	}
}
