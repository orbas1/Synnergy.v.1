package core

import "testing"

func TestBenefitRegistry(t *testing.T) {
	r := NewBenefitRegistry()
	id := r.RegisterBenefit("alice", "housing", 500, "approver")
	if err := r.Claim(id, "alice"); err != nil {
		t.Fatalf("claim: %v", err)
	}
	if err := r.Approve(id, "auditor"); err != nil {
		t.Fatalf("approve: %v", err)
	}
	b, ok := r.GetBenefit(id)
	if !ok || !b.Claimed || b.Status != BenefitStatusApproved {
		t.Fatalf("benefit not approved properly: %#v", b)
	}
	tele := r.Telemetry()
	if tele.Total != 1 || tele.Approved != 1 {
		t.Fatalf("unexpected telemetry: %+v", tele)
	}
}
