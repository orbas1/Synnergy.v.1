package core

import "testing"

func TestBenefitRegistryLifecycle(t *testing.T) {
	reg := NewBenefitRegistry()
	id, err := reg.RegisterBenefit("beneficiary", "housing", 500, "")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := reg.Approve(id, "approver"); err != nil {
		t.Fatalf("approve: %v", err)
	}
	if err := reg.Claim(id, "beneficiary"); err != nil {
		t.Fatalf("claim: %v", err)
	}
	b, ok := reg.GetBenefit(id)
	if !ok || b.Status != BenefitStatusClaimed {
		t.Fatalf("unexpected benefit: %#v", b)
	}
	snap := reg.Snapshot()
	restored := NewBenefitRegistry()
	restored.Restore(snap)
	b2, ok := restored.GetBenefit(id)
	if !ok || b2.Status != BenefitStatusClaimed {
		t.Fatalf("restore mismatch: %#v", b2)
	}
}

func TestBenefitRegistryStatus(t *testing.T) {
	reg := NewBenefitRegistry()
	id1, _ := reg.RegisterBenefit("beneficiary", "housing", 500, "")
	id2, _ := reg.RegisterBenefit("beneficiary", "food", 200, "approver")
	_ = reg.Approve(id1, "approver")
	_ = reg.Claim(id2, "approver")
	status := reg.StatusSummary()
	if status.Total != 2 || status.Approved != 2 || status.Claimed != 1 {
		t.Fatalf("unexpected status: %+v", status)
	}
}
