package core

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"

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

func TestWatchtowerStartStop(t *testing.T) {
	w := NewWatchtowerNode("wt1", nil)
	ctx := context.Background()
	if err := w.Start(ctx); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if !w.running {
		t.Fatalf("watchtower should be running after start")
	}
	if err := w.Start(ctx); err == nil {
		t.Fatalf("expected error on second start")
	}
	if err := w.Stop(); err != nil {
		t.Fatalf("stop failed: %v", err)
	}
	if w.running {
		t.Fatalf("watchtower should not be running after stop")
	}
	// stopping again should not error
	if err := w.Stop(); err != nil {
		t.Fatalf("stop second failed: %v", err)
	}
}

func TestWatchtowerMetrics(t *testing.T) {
	w := NewWatchtowerNode("wt1", nil)
	expected := w.health.Collect(3, 42)
	m := w.Metrics()
	if m.PeerCount != expected.PeerCount || m.LastBlockHeight != expected.LastBlockHeight {
		t.Fatalf("metrics mismatch: got %+v want %+v", m, expected)
	}
	if m.Timestamp.IsZero() {
		t.Fatalf("expected timestamp to be set")
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
}

// compile-time interface check
var _ watchtower.WatchtowerNode = (*Watchtower)(nil)
