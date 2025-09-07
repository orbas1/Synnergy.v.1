package cli

import (
	"strings"
	"testing"

	"synnergy/internal/tokens"
)

// TestSyn845Workflow covers debt token issuance and payments.
func TestSyn845Workflow(t *testing.T) {
	debtRegistry = tokens.NewDebtRegistry()

	out, err := execCommand("syn845", "create", "--name", "Debt", "--symbol", "DBT", "--owner", "alice", "--supply", "1000")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !strings.Contains(out, "token created") {
		t.Fatalf("unexpected create output: %s", out)
	}

	out, err = execCommand("syn845", "issue", "DEBT-1", "loan1", "bob", "100", "5", "1", "2030-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if !strings.Contains(out, "debt issued") {
		t.Fatalf("unexpected issue output: %s", out)
	}

	out, err = execCommand("syn845", "pay", "DEBT-1", "loan1", "10")
	if err != nil {
		t.Fatalf("pay: %v", err)
	}
	if !strings.Contains(out, "payment recorded") {
		t.Fatalf("unexpected pay output: %s", out)
	}

	out, err = execCommand("syn845", "info", "DEBT-1", "loan1")
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	if !strings.Contains(out, "Borrower:bob") || !strings.Contains(out, "Paid:10") {
		t.Fatalf("unexpected info output: %s", out)
	}
}
