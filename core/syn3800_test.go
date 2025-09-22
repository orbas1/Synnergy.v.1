package core

import "testing"

func TestGrantRegistryLifecycle(t *testing.T) {
	r := NewGrantRegistry()
	id := r.CreateGrant("bob", "research", 100)
	if err := r.Authorize(id, "auth1"); err != nil {
		t.Fatalf("authorize: %v", err)
	}
	if err := r.Disburse(id, 40, "phase1", "auth1"); err != nil {
		t.Fatalf("disburse: %v", err)
	}
	g, ok := r.GetGrant(id)
	if !ok {
		t.Fatalf("grant not found")
	}
	if g.Released != 40 || g.Status != GrantStatusActive {
		t.Fatalf("unexpected grant state: %#v", g)
	}
	if err := r.Disburse(id, 60, "phase2", "auth1"); err != nil {
		t.Fatalf("disburse2: %v", err)
	}
	events, ok := r.Audit(id)
	if !ok || len(events) < 3 {
		t.Fatalf("expected audit events")
	}
	tele := r.Telemetry()
	if tele.Total != 1 || tele.Completed != 1 {
		t.Fatalf("unexpected telemetry: %#v", tele)
	}
	list := r.ListGrants()
	if len(list) != 1 || list[0].Status != GrantStatusCompleted {
		t.Fatalf("expected 1 completed grant, got %#v", list)
	}
}
