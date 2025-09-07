package cli

import (
	"bytes"
	"strings"
	"testing"

	"synnergy/internal/tokens"
)

// TestSyn200RegisterMissingFlags verifies that missing required flags return an error.
func TestSyn200RegisterMissingFlags(t *testing.T) {
	carbonRegistry = tokens.NewCarbonRegistry()

	cmd := RootCmd()
	cmd.SetArgs([]string{"syn200", "register"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing required flags")
	}
}

// TestSyn200RegisterWorkflow ensures basic register and info flows work through the CLI.
func TestSyn200RegisterWorkflow(t *testing.T) {
	carbonRegistry = tokens.NewCarbonRegistry()

	cmd := RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	cmd.SetArgs([]string{"syn200", "register", "--owner", "alice", "--name", "proj", "--total", "100"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("register command failed: %v", err)
	}
	id := strings.TrimSpace(buf.String())
	if id == "" {
		t.Fatalf("expected project ID output")
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn200", "info", id})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("info command failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Owner:alice") {
		t.Fatalf("unexpected info output: %s", out)
	}
}
