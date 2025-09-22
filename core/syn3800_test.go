package core

import "testing"

func TestGrantRegistryLifecycle(t *testing.T) {
	r := NewGrantRegistry()
	id := r.CreateGrant("bob", "research", 100)
	if summary := r.Summary(); summary.Total != 1 || summary.Pending != 1 {
		t.Fatalf("unexpected summary %+v", summary)
	}
	if err := r.AddAuthorizer(id, "auth1"); err != nil {
		t.Fatalf("authorizer: %v", err)
	}
	if err := r.DisburseWithActor(id, 40, "phase1", "auth1"); err != nil {
		t.Fatalf("disburse: %v", err)
	}
	g, ok := r.GetGrant(id)
	if !ok || g.Released != 40 || g.Status != GrantStatusActive {
		t.Fatalf("unexpected grant state: %#v", g)
	}
	if err := r.DisburseWithActor(id, 60, "phase2", "auth1"); err != nil {
		t.Fatalf("disburse final: %v", err)
	}
	g, _ = r.GetGrant(id)
	if g.Status != GrantStatusCompleted {
		t.Fatalf("expected completed status")
	}
	events, ok := r.AuditTrail(id)
	if !ok || len(events) < 3 {
		t.Fatalf("expected audit events")
	}
	list := r.ListGrants()
	if len(list) != 1 {
		t.Fatalf("expected 1 grant, got %d", len(list))
	}
	summary := r.Summary()
	if summary.Completed != 1 || summary.Active != 0 || summary.Pending != 0 {
		t.Fatalf("unexpected summary %+v", summary)
	}
}

func TestGrantRegistryAuthorization(t *testing.T) {
	r := NewGrantRegistry()
	id := r.CreateGrant("bob", "research", 10)
	if err := r.DisburseWithActor(id, 5, "", ""); err != nil {
		t.Fatalf("expected disburse without authorizer: %v", err)
	}
	if err := r.AddAuthorizer(id, "wallet1"); err != nil {
		t.Fatalf("authorizer: %v", err)
	}
	if err := r.DisburseWithActor(id, 1, "", "wallet2"); err == nil {
		t.Fatalf("expected unauthorized error")
	}
}
