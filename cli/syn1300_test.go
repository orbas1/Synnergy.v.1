package cli

import (
	"bytes"
	"strings"
	"testing"

	"synnergy/core"
)

// TestSyn1300RegisterRequiresFlags ensures register fails without mandatory flags.
func TestSyn1300RegisterRequiresFlags(t *testing.T) {
	syn1300 = core.NewSupplyChainRegistry()
	cmd := RootCmd()
	cmd.SetArgs([]string{"syn1300", "register"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing required flags")
	}
}

// TestSyn1300Workflow covers register, update and get commands.
func TestSyn1300Workflow(t *testing.T) {
	syn1300 = core.NewSupplyChainRegistry()
	cmd := RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	cmd.SetArgs([]string{"syn1300", "register", "--id", "A1", "--desc", "widget", "--owner", "alice", "--loc", "factory"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn1300", "update", "A1", "--loc", "warehouse", "--status", "shipped"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn1300", "get", "A1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("get failed: %v", err)
	}
	expected := "A1 owned by alice at warehouse status shipped events 2"
	if strings.TrimSpace(buf.String()) != expected {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}
