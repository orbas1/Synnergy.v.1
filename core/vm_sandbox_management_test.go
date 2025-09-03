package core

import (
	"testing"
	"time"
)

func TestSandboxManager(t *testing.T) {
	m := NewSandboxManager(time.Millisecond)
	sb, err := m.StartSandbox("sb1", "addr", 10, 64)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if !sb.Active {
		t.Fatalf("expected active sandbox")
	}
	if err := m.ResetSandbox("sb1"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if err := m.StopSandbox("sb1"); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if sb, ok := m.SandboxStatus("sb1"); !ok || sb.Active {
		t.Fatalf("status mismatch")
	}
	if err := m.DeleteSandbox("sb1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, ok := m.SandboxStatus("sb1"); ok {
		t.Fatalf("sandbox should be removed")
	}
	// verify purge removes stopped sandboxes past TTL
	sb, _ = m.StartSandbox("sb2", "addr", 1, 1)
	_ = m.StopSandbox("sb2")
	time.Sleep(2 * time.Millisecond)
	m.PurgeInactive()
	if _, ok := m.SandboxStatus("sb2"); ok {
		t.Fatalf("expected sb2 purged")
	}
}
