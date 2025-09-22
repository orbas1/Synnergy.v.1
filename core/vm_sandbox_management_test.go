package core

import (
	"errors"
	"testing"
	"time"
)

func TestSandboxManager(t *testing.T) {
	m := NewSandboxManager(2 * time.Millisecond)
	events := make(chan SandboxEvent, 32)
	m.RegisterWatcher(func(ev SandboxEvent) { events <- ev })

	sb, err := m.StartSandbox("sb1", "addr", 10, 64)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if !sb.Active {
		t.Fatalf("expected active sandbox")
	}

	if err := m.Heartbeat("sb1"); err != nil {
		t.Fatalf("heartbeat: %v", err)
	}

	failure := errors.New("execution reverted")
	if err := m.RecordFailure("sb1", failure); err != nil {
		t.Fatalf("record failure: %v", err)
	}
	if sb.Active {
		t.Fatalf("sandbox should be inactive after failure")
	}
	if sb.LastError == "" {
		t.Fatalf("expected failure to be recorded")
	}

	if err := m.RestartSandbox("sb1"); err != nil {
		t.Fatalf("restart: %v", err)
	}
	if !sb.Active || sb.RestartCount != 1 {
		t.Fatalf("sandbox not restarted correctly")
	}

	if err := m.ResetSandbox("sb1"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if err := m.StopSandbox("sb1"); err != nil {
		t.Fatalf("stop: %v", err)
	}

	// allow TTL expiry and purge
	time.Sleep(3 * time.Millisecond)
	m.PurgeInactive()
	if _, ok := m.SandboxStatus("sb1"); ok {
		t.Fatalf("expected sandbox purged")
	}

	metrics := m.Metrics()
	if metrics.TotalCreated != 1 {
		t.Fatalf("expected total created 1, got %d", metrics.TotalCreated)
	}
	if metrics.Failures != 1 {
		t.Fatalf("expected failures 1, got %d", metrics.Failures)
	}
	if metrics.Restarts != 1 {
		t.Fatalf("expected restarts 1, got %d", metrics.Restarts)
	}
	if metrics.Purged != 1 {
		t.Fatalf("expected purged 1, got %d", metrics.Purged)
	}
	if metrics.Active != 0 {
		t.Fatalf("expected no active sandboxes, got %d", metrics.Active)
	}

	close(events)
	seen := map[SandboxEventType]int{}
	for ev := range events {
		seen[ev.Type]++
	}
	for _, typ := range []SandboxEventType{SandboxEventStarted, SandboxEventFailure, SandboxEventRestarted, SandboxEventPurged} {
		if seen[typ] == 0 {
			t.Fatalf("expected event type %s to be emitted", typ)
		}
	}
}
