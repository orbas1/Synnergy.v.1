package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"synnergy/core"
)

// execLoanCLI runs the CLI with args and captures stdout.
func execLoanCLI(args ...string) (string, error) {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetErr(new(bytes.Buffer))
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	w.Close()
	os.Stdout = stdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

// TestLoanProposalCLI exercises creation, voting and retrieval.
func TestLoanProposalCLI(t *testing.T) {
	proposals = make(map[uint64]*core.LoanProposal)
	nextProposalID = 0

	out, err := execLoanCLI("loanproposal", "new", "alice", "bob", "short", "100", "desc", "24")
	if err != nil {
		t.Fatalf("new failed: %v", err)
	}
	if !strings.Contains(out, "created") {
		t.Fatalf("unexpected output: %s", out)
	}

	if _, ok := proposals[1]; !ok {
		t.Fatalf("proposal not stored")
	}

	if _, err := execLoanCLI("loanproposal", "vote", "1", "carol"); err != nil {
		t.Fatalf("vote failed: %v", err)
	}

	out, err = execLoanCLI("--json", "loanproposal", "votes", "1")
	if err != nil {
		t.Fatalf("votes failed: %v", err)
	}
	if !strings.Contains(out, "\"votes\": 1") {
		t.Fatalf("unexpected votes output: %s", out)
	}

	out, err = execLoanCLI("--json", "loanproposal", "get", "1")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if !strings.Contains(out, "alice") {
		t.Fatalf("unexpected get output: %s", out)
	}
}
