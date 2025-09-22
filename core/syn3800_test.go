package core

import "testing"

func TestGrantRegistryLifecycle(t *testing.T) {
	reg := NewGrantRegistry()
	id, err := reg.CreateGrant("beneficiary", "research", 100, "")
	if err != nil {
		t.Fatalf("create grant: %v", err)
	}
	if err := reg.Authorize(id, "auditor"); err != nil {
		t.Fatalf("authorize: %v", err)
	}
	if err := reg.Disburse(id, 40, "phase1", "auditor"); err != nil {
		t.Fatalf("disburse: %v", err)
	}
	g, ok := reg.GetGrant(id)
	if !ok || g.Released != 40 || g.Status != GrantStatusAuthorized {
		t.Fatalf("unexpected grant state: %#v", g)
	}
	events := reg.AuditLog(id)
	if len(events) < 2 {
		t.Fatalf("expected audit events")
	}
	snap := reg.Snapshot()
	restored := NewGrantRegistry()
	restored.Restore(snap)
	g2, ok := restored.GetGrant(id)
	if !ok || g2.Released != 40 {
		t.Fatalf("restore mismatch: %#v", g2)
	}
}

func TestGrantRegistryStatus(t *testing.T) {
	reg := NewGrantRegistry()
	id1, _ := reg.CreateGrant("beneficiary", "research", 100, "")
	id2, _ := reg.CreateGrant("beneficiary", "research", 200, "auditor")
	_ = reg.Authorize(id1, "auditor")
	_ = reg.Disburse(id2, 200, "full", "auditor")
	status := reg.StatusSummary()
	if status.Total != 2 || status.Active != 1 || status.Completed != 1 {
		t.Fatalf("unexpected status: %+v", status)
	}
}
