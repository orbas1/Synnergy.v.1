package core

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
	"time"

	watchtower "synnergy/internal/nodes/extra/watchtower"
)

func TestNewWatchtowerNode(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	w := NewWatchtowerNode("wt1", logger)
	if w.ID() != "wt1" {
		t.Fatalf("expected id wt1 got %s", w.ID())
	}
	if w.firewall == nil {
		t.Fatalf("expected firewall initialised")
	}
	if w.health == nil {
		t.Fatalf("expected health logger initialised")
	}
	if w.logger != logger {
		t.Fatalf("logger mismatch")
	}
	if w.Firewall() != w.firewall {
		t.Fatalf("firewall accessor mismatch")
	}
}

func TestWatchtowerStartStopEmitsEvents(t *testing.T) {
	w := NewWatchtowerNode("wt1", nil)
	w.tickInterval = 100 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := w.Start(ctx); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	time.Sleep(250 * time.Millisecond)
	if err := w.Stop(); err != nil {
		t.Fatalf("stop failed: %v", err)
	}
	events := w.Events()
	var hasStart, hasStop bool
	for _, ev := range events {
		if ev.Type == WatchtowerEventStarted {
			hasStart = true
		}
		if ev.Type == WatchtowerEventStopped {
			hasStop = true
		}
	}
	if !hasStart || !hasStop {
		t.Fatalf("expected start and stop events, got %+v", events)
	}
}

func TestWatchtowerRunIntegritySweep(t *testing.T) {
	w := NewWatchtowerNode("wt", nil)
	node := NewNode("n", "addr", NewLedger())
	for i := 0; i < 5; i++ {
		node.Mempool = append(node.Mempool, &Transaction{})
	}
	w.AttachNode(node)
	events, err := w.RunIntegritySweep(context.Background(), 3)
	if err != nil {
		t.Fatalf("integrity sweep: %v", err)
	}
	if len(events) == 0 {
		t.Fatalf("expected alert event")
	}
	found := false
	for _, ev := range events {
		if ev.Type == WatchtowerEventAlert {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected alert event, got %+v", events)
	}
}

func TestWatchtowerReportFork(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	w := NewWatchtowerNode("id", logger)
	w.ReportFork(10, "abc")
	if !strings.Contains(buf.String(), "fork detected at height 10 hash abc") {
		t.Fatalf("unexpected log: %q", buf.String())
	}
	events := w.Events()
	if len(events) == 0 || events[len(events)-1].Type != WatchtowerEventForkDetected {
		t.Fatalf("expected fork event, got %+v", events)
	}
}

// compile-time interface check
var _ watchtower.WatchtowerNode = (*Watchtower)(nil)
