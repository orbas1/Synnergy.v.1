package core

import (
	"errors"
	"testing"
)

func TestGrantRegistryLifecycle(t *testing.T) {
	r := NewGrantRegistry()
	if _, err := r.CreateGrant("", "research", 10); err == nil {
		t.Fatal("expected validation error for empty beneficiary")
	}
	if _, err := r.CreateGrant("bob", "", 10); err == nil {
		t.Fatal("expected validation error for empty name")
	}
	if _, err := r.CreateGrant("bob", "research", 0); err == nil {
		t.Fatal("expected validation error for zero amount")
	}

	id, err := r.CreateGrant("bob", "research", 100, "auth1")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !r.IsAuthorized(id, "auth1") {
		t.Fatalf("expected initial authorizer registered")
	}
	if _, err := r.Authorize(id, ""); err == nil {
		t.Fatal("expected error for empty address")
	}
	if _, err := r.Authorize(id, "auth2"); err != nil {
		t.Fatalf("authorize: %v", err)
	}
	if _, err := r.Disburse(id, 0, "", "auth1"); !errors.Is(err, ErrGrantInvalidAmount) {
		t.Fatalf("expected invalid amount, got %v", err)
	}
	if _, err := r.Disburse(id, 10, "phase1", "bad"); !errors.Is(err, ErrGrantUnauthorized) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
	evt, err := r.Disburse(id, 40, "phase1", "auth1")
	if err != nil {
		t.Fatalf("disburse: %v", err)
	}
	if evt.Amount != 40 || evt.Signer != "auth1" {
		t.Fatalf("unexpected event: %#v", evt)
	}
	g, ok := r.GetGrant(id)
	if !ok {
		t.Fatalf("grant not found")
	}
	if g.Status != GrantStatusActive || g.Released != 40 {
		t.Fatalf("unexpected grant state: %#v", g)
	}
	summary := r.StatusSummary()
	if summary.Total != 1 || summary.Active != 1 || summary.Pending != 0 {
		t.Fatalf("unexpected summary after partial release: %#v", summary)
	}
	if _, err := r.Disburse(id, 60, "phase2", "auth2"); err != nil {
		t.Fatalf("disburse second: %v", err)
	}
	g, _ = r.GetGrant(id)
	if g.Status != GrantStatusCompleted || g.Released != 100 {
		t.Fatalf("expected completed grant, got %#v", g)
	}
	events, err := r.Audit(id)
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	if len(events) < 3 {
		t.Fatalf("expected at least 3 events, got %d", len(events))
	}
	summary = r.StatusSummary()
	if summary.Completed != 1 || summary.Active != 0 || summary.Pending != 0 {
		t.Fatalf("unexpected summary after completion: %#v", summary)
	}

	snap := r.Snapshot()
	restored := NewGrantRegistry()
	restored.Restore(snap)
	g2, ok := restored.GetGrant(id)
	if !ok || g2.Status != GrantStatusCompleted || g2.Released != 100 {
		t.Fatalf("restored grant mismatch: %#v", g2)
	}
	if _, err := restored.Disburse(id, 1, "extra", "auth1"); !errors.Is(err, ErrGrantInvalidAmount) {
		t.Fatalf("expected invalid amount on over release, got %v", err)
	}
}
