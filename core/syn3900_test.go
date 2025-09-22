package core

import "testing"

func TestBenefitRegistryLifecycle(t *testing.T) {
	r := NewBenefitRegistry()
	id := r.RegisterBenefit("recipient", "housing", 200)
	if err := r.Claim(id, "recipient"); err != nil {
		t.Fatalf("claim: %v", err)
	}
	if err := r.Approve(id, "approver"); err != nil {
		t.Fatalf("approve: %v", err)
	}
	b, ok := r.GetBenefit(id)
	if !ok {
		t.Fatalf("benefit not found")
	}
	if !b.Claimed || b.Status != BenefitStatusApproved {
		t.Fatalf("unexpected benefit state: %#v", b)
	}
	tele := r.Telemetry()
	if tele.Total != 1 || tele.Approved != 1 || tele.Claimed != 1 {
		t.Fatalf("unexpected telemetry: %#v", tele)
	}
	list := r.ListBenefits()
	if len(list) != 1 || list[0].Program != "housing" {
		t.Fatalf("unexpected list: %#v", list)
	}
}
