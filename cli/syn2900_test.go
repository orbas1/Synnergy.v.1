package cli

import (
	"testing"
	"time"

	"synnergy/internal/tokens"
)

func run2900(args ...string) error {
	cmd := RootCmd()
	cmd.SetArgs(args)
	return cmd.Execute()
}

func TestSyn2900IssueAndClaim(t *testing.T) {
	insuranceRegistry = tokens.NewInsuranceRegistry()
	start := time.Now().Format(time.RFC3339)
	end := time.Now().Add(time.Hour).Format(time.RFC3339)
	if err := run2900("syn2900", "issue", "--holder", "alice", "--coverage", "basic", "--premium", "10", "--payout", "100", "--start", start, "--end", end); err != nil {
		t.Fatalf("issue failed: %v", err)
	}
	policies := insuranceRegistry.ListPolicies()
	if len(policies) != 1 {
		t.Fatalf("policy not issued")
	}
	id := policies[0].PolicyID
	if err := run2900("syn2900", "claim", id, "fire", "20"); err != nil {
		t.Fatalf("claim failed: %v", err)
	}
}

func TestSyn2900IssueMissingFields(t *testing.T) {
	insuranceRegistry = tokens.NewInsuranceRegistry()
	if err := run2900("syn2900", "issue", "--holder", "", "--coverage", "", "--premium", "0", "--payout", "0"); err == nil {
		t.Fatalf("expected error for missing fields")
	}
}
