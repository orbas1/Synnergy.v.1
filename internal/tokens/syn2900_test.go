package tokens

import (
	"testing"
	"time"
)

func TestInsuranceRegistryLifecycle(t *testing.T) {
	registry := NewInsuranceRegistry()
	start := time.Now().Add(time.Hour)
	end := start.Add(72 * time.Hour)

	policy, err := registry.IssuePolicy("alice", "property", 250, 5_000, 100, 10_000, start, end)
	if err != nil {
		t.Fatalf("issue policy: %v", err)
	}
	if policy.PolicyID == "" {
		t.Fatalf("expected policy id")
	}

	if err := registry.PayPremium(policy.PolicyID, 0); err != ErrInsurancePremiumInvalid {
		t.Fatalf("expected invalid premium error, got %v", err)
	}
	if err := registry.PayPremium(policy.PolicyID, 500); err != nil {
		t.Fatalf("pay premium: %v", err)
	}
	clone, ok := registry.GetPolicy(policy.PolicyID)
	if !ok {
		t.Fatalf("policy not found")
	}
	if clone.PaidPremium != 500 {
		t.Fatalf("expected paid premium 500, got %d", clone.PaidPremium)
	}

	if _, err := registry.FileClaim(policy.PolicyID, "fire", 20_000); err != ErrInsuranceLimitExceeded {
		t.Fatalf("expected limit exceeded error, got %v", err)
	}
	claim, err := registry.FileClaim(policy.PolicyID, "fire", 4_000)
	if err != nil {
		t.Fatalf("file claim: %v", err)
	}

	exposure := registry.TotalExposure(time.Now())
	if exposure != 6_000 {
		t.Fatalf("expected exposure 6000, got %d", exposure)
	}

	if err := registry.SettleClaim(policy.PolicyID, claim.ClaimID); err != nil {
		t.Fatalf("settle claim: %v", err)
	}

	claim2, err := registry.FileClaim(policy.PolicyID, "water", 6_000)
	if err != nil {
		t.Fatalf("second claim: %v", err)
	}
	if err := registry.SettleClaim(policy.PolicyID, claim2.ClaimID); err != nil {
		t.Fatalf("settle second claim: %v", err)
	}
	active := registry.ListActivePolicies(time.Now())
	if len(active) != 0 {
		t.Fatalf("policy should be inactive after settling limit")
	}

	if err := registry.Deactivate(policy.PolicyID); err != nil {
		t.Fatalf("deactivate: %v", err)
	}
	if err := registry.SettleClaim(policy.PolicyID, "unknown"); err != ErrInsuranceClaimNotFound {
		t.Fatalf("expected claim not found error, got %v", err)
	}
}
