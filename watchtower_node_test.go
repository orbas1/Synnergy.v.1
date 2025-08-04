package synnergy

import (
	"context"
	"io"
	"log"
	"testing"
)

func TestWatchtowerNode(t *testing.T) {
	w := NewWatchtowerNode("wt1", log.New(io.Discard, "", 0))
	if w.ID() != "wt1" {
		t.Fatalf("unexpected id: %s", w.ID())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := w.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := w.Start(ctx); err == nil {
		t.Fatalf("expected error on double start")
	}

	w.health.Collect(1, 1)
	if w.Metrics().PeerCount != 1 {
		t.Fatalf("metrics not recorded")
	}
	w.ReportFork(5, "hash")

	if err := w.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
