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
	if sb.Active {
		t.Fatalf("sandbox should be inactive")
	}
	if err := m.DeleteSandbox("sb1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, ok := m.SandboxStatus("sb1"); ok {
		t.Fatalf("sandbox should be removed")
	}
}
