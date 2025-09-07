package cli

import (
	"testing"
	"time"
)

func run2700(args ...string) error {
	cmd := RootCmd()
	cmd.SetArgs(args)
	return cmd.Execute()
}

func TestSyn2700CreateAndClaim(t *testing.T) {
	schedule = nil
	entry := time.Now().Add(-time.Hour).Format(time.RFC3339) + "=10"
	if err := run2700("syn2700", "create", "--entries", entry); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := run2700("syn2700", "claim"); err != nil {
		t.Fatalf("claim failed: %v", err)
	}
}

func TestSyn2700CreateMissingEntries(t *testing.T) {
	schedule = nil
	if err := run2700("syn2700", "create", "--entries", ""); err == nil {
		t.Fatalf("expected error for missing entries")
	}
}
