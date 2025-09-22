package tokens

import (
	"testing"
	"time"
)

func TestLifePolicyRegistryLifecycle(t *testing.T) {
	registry := NewLifePolicyRegistry()
	start := time.Now().Add(2 * time.Hour)
	end := start.Add(48 * time.Hour)

	policy, err := registry.IssuePolicy("alice", "bob", 1_000, 100, start, end)
	if err != nil {
		t.Fatalf("issue policy: %v", err)
	}
	if policy.PolicyID == "" {
		t.Fatalf("expected policy id")
	}

	if err := registry.PayPremium(policy.PolicyID, 0); err != ErrInvalidPremium {
		t.Fatalf("expected ErrInvalidPremium, got %v", err)
	}
	if err := registry.PayPremium(policy.PolicyID, 200); err != nil {
		t.Fatalf("pay premium: %v", err)
	}

	clone, ok := registry.GetPolicy(policy.PolicyID)
	if !ok {
		t.Fatalf("policy not retrievable")
	}
	if clone.PaidPremium != 200 {
		t.Fatalf("expected paid premium 200, got %d", clone.PaidPremium)
	}

	if _, err := registry.FileClaim(policy.PolicyID, 2_000); err != ErrCoverageExceeded {
		t.Fatalf("expected coverage exceeded, got %v", err)
	}
	claim, err := registry.FileClaim(policy.PolicyID, 400)
	if err != nil {
		t.Fatalf("file claim: %v", err)
	}

	if err := registry.SettleClaim(policy.PolicyID, claim.ClaimID); err != nil {
		t.Fatalf("settle claim: %v", err)
	}

	// Settling the remaining coverage should deactivate the policy.
	claim2, err := registry.FileClaim(policy.PolicyID, 600)
	if err != nil {
		t.Fatalf("second claim: %v", err)
	}
	if err := registry.SettleClaim(policy.PolicyID, claim2.ClaimID); err != nil {
		t.Fatalf("settle second claim: %v", err)
	}
	activePolicies := registry.ListActivePolicies(time.Now())
	if len(activePolicies) != 0 {
		t.Fatalf("expected policy to be inactive: %+v", activePolicies)
	}

	if err := registry.Deactivate(policy.PolicyID); err != nil {
		t.Fatalf("deactivate policy: %v", err)
	}
	if err := registry.SettleClaim(policy.PolicyID, "unknown"); err != ErrClaimNotFound {
		t.Fatalf("expected ErrClaimNotFound, got %v", err)
	}
}

func TestLifePolicyGracePeriod(t *testing.T) {
	registry := NewLifePolicyRegistry()
	start := time.Now()
	end := start.Add(24 * time.Hour)
	registry.SetDefaultGrace(4 * time.Hour)
	policy, err := registry.IssuePolicy("alice", "bob", 500, 50, start, end)
	if err != nil {
		t.Fatalf("issue policy: %v", err)
	}

	// No premium paid yet, still in grace period for 4 hours.
	if !policy.InGracePeriod(start.Add(3 * time.Hour)) {
		t.Fatalf("expected in grace period")
	}

	if err := registry.PayPremium(policy.PolicyID, 10); err != nil {
		t.Fatalf("pay premium: %v", err)
	}
	updated, _ := registry.GetPolicy(policy.PolicyID)
	if updated.LastPremium.Before(start) {
		t.Fatalf("expected last premium to update")
	}
}
