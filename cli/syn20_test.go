package cli

import (
	"bytes"
	"strings"
	"testing"

	"synnergy/internal/tokens"
)

// TestSyn20InitRequiresFlags ensures init fails when required flags are missing.
func TestSyn20InitRequiresFlags(t *testing.T) {
	tokenRegistry = tokens.NewRegistry()
	syn20 = nil

	cmd := RootCmd()
	cmd.SetArgs([]string{"syn20", "init"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing required flags")
	}
}

// TestSyn20MintWorkflow verifies init, mint and balance subcommands.
func TestSyn20MintWorkflow(t *testing.T) {
	tokenRegistry = tokens.NewRegistry()
	syn20 = nil

	cmd := RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	cmd.SetArgs([]string{"syn20", "init", "--name", "Utility", "--symbol", "UTL"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn20", "mint", "alice", "100"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("mint command failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn20", "balance", "alice"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("balance command failed: %v", err)
	}
	if strings.TrimSpace(buf.String()) != "100" {
		t.Fatalf("unexpected balance output: %s", buf.String())
	}
}
