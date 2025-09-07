package cli

import (
	"bytes"
	"strings"
	"testing"

	"synnergy/internal/tokens"
)

// TestSyn1100AddRequiresFlags ensures add fails when flags are missing.
func TestSyn1100AddRequiresFlags(t *testing.T) {
	syn1100 = tokens.NewSYN1100Token()

	cmd := RootCmd()
	cmd.SetArgs([]string{"syn1100", "add", "--id", "1", "--owner", "alice"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing data flag")
	}
}

// TestSyn1100Workflow verifies record access control operations.
func TestSyn1100Workflow(t *testing.T) {
	syn1100 = tokens.NewSYN1100Token()

	cmd := RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	cmd.SetArgs([]string{"syn1100", "add", "--id", "1", "--owner", "alice", "--data", "record"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("add command failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn1100", "grant", "--id", "1", "--grantee", "bob"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("grant command failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn1100", "get", "--id", "1", "--caller", "bob"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("get command failed: %v", err)
	}
	if strings.TrimSpace(buf.String()) != "record" {
		t.Fatalf("unexpected record output: %s", buf.String())
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn1100", "revoke", "--id", "1", "--grantee", "bob"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("revoke command failed: %v", err)
	}

	cmd.SetArgs([]string{"syn1100", "get", "--id", "1", "--caller", "bob"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error after access revoked")
	}
}
