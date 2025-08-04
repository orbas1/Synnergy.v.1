package core

import "testing"

func TestGrantRegistry(t *testing.T) {
	r := NewGrantRegistry()
	id := r.CreateGrant("bob", "research", 100)
	if err := r.Disburse(id, 40, "phase1"); err != nil {
		t.Fatalf("disburse: %v", err)
	}
	g, ok := r.GetGrant(id)
	if !ok || g.Released != 40 {
		t.Fatalf("unexpected grant state: %#v", g)
	}
	list := r.ListGrants()
	if len(list) != 1 {
		t.Fatalf("expected 1 grant, got %d", len(list))
	}
}
