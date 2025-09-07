package cli

import (
	"bytes"
	"strings"
	"testing"
)

// TestSyn1700InitRequiresFlags ensures init fails when required flags are missing.
func TestSyn1700InitRequiresFlags(t *testing.T) {
	event = nil
	cmd := RootCmd()
	cmd.SetArgs([]string{"syn1700", "init", "--name", "Concert", "--desc", "desc", "--location", "NY", "--start", "1", "--end", "2"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing supply flag")
	}
}

// TestSyn1700Workflow covers init, issue and verify operations.
func TestSyn1700Workflow(t *testing.T) {
	event = nil
	cmd := RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	cmd.SetArgs([]string{"syn1700", "init", "--name", "Concert", "--desc", "desc", "--location", "NY", "--start", "1", "--end", "2", "--supply", "100"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn1700", "issue", "bob", "A", "VIP", "50"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("issue failed: %v", err)
	}
	ticketID := strings.TrimSpace(buf.String())

	buf.Reset()
	cmd.SetArgs([]string{"syn1700", "verify", ticketID, "bob"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if strings.TrimSpace(buf.String()) != "true" {
		t.Fatalf("expected true verification, got %s", buf.String())
	}
}
