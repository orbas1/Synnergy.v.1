package watchtower

import (
	"context"
	"testing"
)

func TestBasicWatchtowerLifecycle(t *testing.T) {
	w := NewBasicWatchtower("w1")
	if err := w.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	w.ReportFork(1, "hash")
	m := w.Metrics()
	if m.Timestamp.IsZero() {
		t.Fatalf("metrics timestamp not set")
	}
	if err := w.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
