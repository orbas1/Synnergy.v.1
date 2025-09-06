package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"synnergy/core"
)

// execLoanpoolCLI executes rootCmd with args capturing stdout.
func execLoanpoolCLI(args ...string) (string, error) {
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

// TestLoanpoolCLIFlow covers submit, vote and get operations.
func TestLoanpoolCLIFlow(t *testing.T) {
	loanPool = core.NewLoanPool(1_000_000)

	out, err := execLoanpoolCLI("loanpool", "submit", "alice", "bob", "short", "100", "desc")
	if err != nil {
		t.Fatalf("submit failed: %v", err)
	}
	if !strings.Contains(out, "submitted") {
		t.Fatalf("unexpected submit output: %s", out)
	}

	if _, err := execLoanpoolCLI("loanpool", "vote", "alice", "1"); err != nil {
		t.Fatalf("vote failed: %v", err)
	}

	out, err = execLoanpoolCLI("--json", "loanpool", "get", "1")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if !strings.Contains(out, "alice") {
		t.Fatalf("unexpected get output: %s", out)
	}
}
