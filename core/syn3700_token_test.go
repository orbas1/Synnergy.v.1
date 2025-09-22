package core

import "testing"

func TestSYN3700TokenLifecycle(t *testing.T) {
	tok := NewSYN3700Token("Index", "IDX")
	tok.AddController("ctrl1")
	if !tok.HasController("ctrl1") {
		t.Fatalf("expected controller registered")
	}
	if err := tok.AddComponent("ctrl1", "AAA", 0.5, 0.1); err != nil {
		t.Fatalf("add component: %v", err)
	}
	if err := tok.AddComponent("ctrl1", "BBB", 0.5, 0.2); err != nil {
		t.Fatalf("add component: %v", err)
	}
	snap := tok.Snapshot()
	if snap.Symbol != "IDX" || len(snap.Components) != 2 {
		t.Fatalf("unexpected snapshot: %+v", snap)
	}
	tele := tok.Telemetry()
	if tele.ComponentCount != 2 || tele.ControllerCount != 1 {
		t.Fatalf("unexpected telemetry: %+v", tele)
	}
	updates := tok.Rebalance("ctrl1")
	if len(updates) != 2 {
		t.Fatalf("expected rebalance updates")
	}
	audit := tok.AuditTrail()
	if len(audit) < 3 {
		t.Fatalf("expected audit entries")
	}
	if err := tok.RemoveComponent("ctrl1", "BBB"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	comps := tok.ListComponents()
	if len(comps) != 1 || comps[0].Token != "AAA" {
		t.Fatalf("unexpected components %+v", comps)
	}
	val := tok.Value(map[string]float64{"AAA": 2})
	if val != 1 {
		t.Fatalf("expected value 1 got %f", val)
	}
}
