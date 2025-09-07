package cli

import (
	"testing"
	"time"

	"synnergy/internal/tokens"
)

func run2800(args ...string) error {
	cmd := RootCmd()
	cmd.SetArgs(args)
	return cmd.Execute()
}

func TestSyn2800IssueAndPay(t *testing.T) {
	lifeRegistry = tokens.NewLifePolicyRegistry()
	start := time.Now().Format(time.RFC3339)
	end := time.Now().Add(time.Hour).Format(time.RFC3339)
	if err := run2800("syn2800", "issue", "--insured", "alice", "--beneficiary", "bob", "--coverage", "100", "--premium", "10", "--start", start, "--end", end); err != nil {
		t.Fatalf("issue failed: %v", err)
	}
	policies := lifeRegistry.ListPolicies()
	if len(policies) != 1 {
		t.Fatalf("policy not issued")
	}
	id := policies[0].PolicyID
	if err := run2800("syn2800", "pay", id, "5"); err != nil {
		t.Fatalf("pay failed: %v", err)
	}
}

func TestSyn2800IssueMissingFields(t *testing.T) {
	lifeRegistry = tokens.NewLifePolicyRegistry()
	if err := run2800("syn2800", "issue", "--insured", "", "--beneficiary", "", "--coverage", "0", "--premium", "0"); err == nil {
		t.Fatalf("expected error for missing fields")
	}
}
