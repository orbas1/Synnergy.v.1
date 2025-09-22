package core

import "testing"

func TestGrantRegistry(t *testing.T) {
	r := NewGrantRegistry()
	id := r.CreateGrant("bob", "research", 100, "auth")
	if err := r.Authorize(id, "delegate"); err != nil {
		t.Fatalf("authorize: %v", err)
	}
	if err := r.Disburse(id, 40, "phase1", "delegate"); err != nil {
		t.Fatalf("disburse: %v", err)
	}
	if err := r.Disburse(id, 60, "", "auth"); err != nil {
		t.Fatalf("final disburse: %v", err)
	}
	g, ok := r.GetGrant(id)
	if !ok || g.Released != 100 || g.Status != GrantStatusCompleted {
		t.Fatalf("unexpected grant state: %#v", g)
	}
	events, err := r.Audit(id)
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	if len(events) < 3 {
		t.Fatalf("expected audit events, got %d", len(events))
	}
	tele := r.Telemetry()
	if tele.Total != 1 || tele.Completed != 1 {
		t.Fatalf("unexpected telemetry: %+v", tele)
	}
	list := r.ListGrants()
	if len(list) != 1 {
		t.Fatalf("expected 1 grant, got %d", len(list))
	}
}
