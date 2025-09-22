package core

import "testing"

func TestSYN3700ControllerLifecycle(t *testing.T) {
	tok := NewSYN3700Token("Index", "IDX")
	tok.AddController("controller1")
	if err := tok.AddComponent("AAA", 0.5, 0.1, "controller1"); err != nil {
		t.Fatalf("add component: %v", err)
	}
	if err := tok.AddComponent("BBB", 1.5, 0.2, "controller1"); err != nil {
		t.Fatalf("add second component: %v", err)
	}
	if tok.ComponentCount() != 2 {
		t.Fatalf("expected 2 components")
	}
	if tok.ControllerCount() != 1 {
		t.Fatalf("expected controller count 1")
	}
	_, err := tok.Rebalance("controller1")
	if err != nil {
		t.Fatalf("rebalance: %v", err)
	}
	audit := tok.Audit()
	if len(audit) == 0 {
		t.Fatalf("expected audit entries")
	}
}
