package core

import "testing"

func TestBenefitRegistryLifecycle(t *testing.T) {
	r := NewBenefitRegistry()
	id := r.RegisterBenefit("alice", "housing", 500)
	if err := r.AddApprover(id, "approver1"); err != nil {
		t.Fatalf("approver: %v", err)
	}
	if err := r.Claim(id, "alice"); err != nil {
		t.Fatalf("claim: %v", err)
	}
	if err := r.Approve(id, "approver1"); err != nil {
		t.Fatalf("approve: %v", err)
	}
	b, ok := r.GetBenefit(id)
	if !ok || !b.Claimed || b.Status != BenefitStatusApproved {
		t.Fatalf("benefit not approved properly: %#v", b)
	}
	if summary := r.Summary(); summary.Approved != 1 || summary.Total != 1 {
		t.Fatalf("unexpected summary %+v", summary)
	}
	events, ok := r.Events(id)
	if !ok || len(events) < 3 {
		t.Fatalf("expected events recorded")
	}
}

func TestBenefitRegistryAuth(t *testing.T) {
	r := NewBenefitRegistry()
	id := r.RegisterBenefit("bob", "housing", 100)
	if err := r.Claim(id, "charlie"); err == nil {
		t.Fatalf("expected authentication failure")
	}
}
