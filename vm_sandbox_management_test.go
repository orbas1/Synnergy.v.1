package synnergy

import "testing"

func TestSandboxManager(t *testing.T) {
	m := NewSandboxManager()
	sb, err := m.StartSandbox("sb1", "c1", 10, 1024)
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
	status, ok := m.SandboxStatus("sb1")
	if !ok {
		t.Fatalf("sandbox should exist")
	}
	if status.Active {
		t.Fatalf("sandbox should be inactive")
	}
	status.Active = true
	fresh, ok := m.SandboxStatus("sb1")
	if !ok || fresh.Active {
		t.Fatalf("sandbox mutation should not leak")
	}
	listed := m.ListSandboxes()
	if len(listed) != 1 {
		t.Fatalf("expected one sandbox, got %d", len(listed))
	}
	listed[0].Active = true
	fresh, _ = m.SandboxStatus("sb1")
	if fresh.Active {
		t.Fatalf("list mutation should not leak")
	}
	if err := m.DeleteSandbox("sb1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, ok := m.SandboxStatus("sb1"); ok {
		t.Fatalf("sandbox should be removed")
	}
}
