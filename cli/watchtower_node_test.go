package cli

import (
	"context"
	"io"
	"log"
	"testing"

	"synnergy/core"
)

func TestWatchtowerNodeLifecycle(t *testing.T) {
	watchNode = core.NewWatchtowerNode("n1", log.New(io.Discard, "", 0))
	if err := watchNode.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	watchNode.ReportFork(1, "abc")
	_ = watchNode.Metrics()
	if err := watchNode.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
