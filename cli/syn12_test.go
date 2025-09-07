package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"synnergy/internal/tokens"
)

// TestSyn12InitRequiresFlags ensures init fails when required flags are missing.
func TestSyn12InitRequiresFlags(t *testing.T) {
	tokenRegistry = tokens.NewRegistry()
	syn12Token = nil

	cmd := RootCmd()
	cmd.SetArgs([]string{"syn12", "init"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing required flags")
	}
}

// TestSyn12MintWorkflow verifies init, mint and balance operations.
func TestSyn12MintWorkflow(t *testing.T) {
	tokenRegistry = tokens.NewRegistry()
	syn12Token = nil

	cmd := RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	issue := time.Now().Format(time.RFC3339)
	maturity := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	cmd.SetArgs([]string{"syn12", "init", "--name", "TBill", "--symbol", "TB", "--bill", "T123", "--issuer", "Gov", "--face", "1000", "--issue", issue, "--maturity", maturity})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn12", "mint", "alice", "100"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("mint command failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn12", "balance", "alice"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("balance command failed: %v", err)
	}
	if strings.TrimSpace(buf.String()) != "100" {
		t.Fatalf("unexpected balance output: %s", buf.String())
	}
}
