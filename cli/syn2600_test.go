package cli

import (
	"testing"
	"time"

	"synnergy/internal/tokens"
)

func run2600(args ...string) error {
	cmd := RootCmd()
	cmd.SetArgs(args)
	return cmd.Execute()
}

func TestSyn2600IssueAndTransfer(t *testing.T) {
	investorRegistry = tokens.NewInvestorRegistry()
	expiry := time.Now().Add(time.Hour).Format(time.RFC3339)
	if err := run2600("syn2600", "issue", "--asset", "gold", "--owner", "alice", "--shares", "10", "--expiry", expiry); err != nil {
		t.Fatalf("issue failed: %v", err)
	}
	tokens := investorRegistry.List()
	if len(tokens) != 1 || tokens[0].Owner != "alice" {
		t.Fatalf("token not issued correctly")
	}
	id := tokens[0].ID
	if err := run2600("syn2600", "transfer", id, "bob"); err != nil {
		t.Fatalf("transfer failed: %v", err)
	}
	tok, _ := investorRegistry.Get(id)
	if tok.Owner != "bob" {
		t.Fatalf("owner not updated")
	}
}

func TestSyn2600IssueMissingFields(t *testing.T) {
	investorRegistry = tokens.NewInvestorRegistry()
	if err := run2600("syn2600", "issue", "--asset", "", "--owner", "", "--shares", "0"); err == nil {
		t.Fatalf("expected error for missing fields")
	}
}
