package cli

import (
	"bytes"
	"strings"
	"testing"

	"synnergy/core"
)

// TestSyn131CreateRequiresFlags ensures create fails when mandatory flags are missing.
func TestSyn131CreateRequiresFlags(t *testing.T) {
	syn131 = core.NewSYN131Registry()
	cmd := RootCmd()
	cmd.SetArgs([]string{"syn131", "create"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing required flags")
	}
}

// TestSyn131Workflow covers create, valuation update and get commands.
func TestSyn131Workflow(t *testing.T) {
	syn131 = core.NewSYN131Registry()
	cmd := RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	cmd.SetArgs([]string{"syn131", "create", "--id", "IP1", "--name", "Patent", "--symbol", "PAT", "--owner", "alice", "--valuation", "1000"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn131", "value", "IP1", "1500"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("valuation update failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn131", "get", "IP1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("get failed: %v", err)
	}
	expected := "IP1 Patent PAT owner:alice val:1500"
	if strings.TrimSpace(buf.String()) != expected {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}
