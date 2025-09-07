package cli

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
	"time"

	"synnergy/core"
)

// TestSyn1401IssueRequiresFlags ensures issue command validates mandatory flags.
func TestSyn1401IssueRequiresFlags(t *testing.T) {
	investments = core.NewInvestmentRegistry()
	cmd := RootCmd()
	// Missing maturity flag
	cmd.SetArgs([]string{"syn1401", "issue", "--id", "INV1", "--owner", "alice", "--principal", "1000", "--rate", "0.1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing maturity flag")
	}
}

// TestSyn1401Workflow covers issuing and retrieving an investment.
func TestSyn1401Workflow(t *testing.T) {
	investments = core.NewInvestmentRegistry()
	cmd := RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	maturity := time.Now().Add(time.Hour).Unix()
	cmd.SetArgs([]string{"syn1401", "issue", "--id", "INV1", "--owner", "alice", "--principal", "1000", "--rate", "0.1", "--maturity", strconv.FormatInt(maturity, 10)})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("issue failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn1401", "get", "INV1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("get failed: %v", err)
	}
	expected := "INV1 owner:alice principal:1000 accrued:0"
	if strings.TrimSpace(buf.String()) != expected {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}
