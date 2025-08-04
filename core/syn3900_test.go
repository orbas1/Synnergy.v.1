package core

import "testing"

func TestBenefitRegistry(t *testing.T) {
	r := NewBenefitRegistry()
	id := r.RegisterBenefit("alice", "housing", 500)
	if err := r.Claim(id); err != nil {
		t.Fatalf("claim: %v", err)
	}
	b, ok := r.GetBenefit(id)
	if !ok || !b.Claimed {
		t.Fatalf("benefit not claimed properly")
	}
}
