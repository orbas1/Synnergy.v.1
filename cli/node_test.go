package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"synnergy/core"
)

// TestNodeCLICommands exercises basic node subcommands and JSON output.
func TestNodeCLICommands(t *testing.T) {
	// reset global ledger and node
	origLedger := ledger
	ledger = core.NewLedger()
	t.Cleanup(func() { ledger = origLedger })
	currentNode = core.NewNode("node1", "addr", ledger)

	out, err := execNodeCLI("--json", "node", "info")
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
	if !strings.Contains(out, "node1") {
		t.Fatalf("unexpected info output: %s", out)
	}

	if _, err := execNodeCLI("node", "stake", "val1", "5"); err != nil {
		t.Fatalf("stake failed: %v", err)
	}
	if currentNode.Validators.Stake("val1") != 5 {
		t.Fatalf("stake not recorded: got %d", currentNode.Validators.Stake("val1"))
	}

	ledger.Mint("alice", 100)
	if _, err := execNodeCLI("node", "addtx", "alice", "bob", "10", "1", "1"); err != nil {
		t.Fatalf("addtx failed: %v", err)
	}
	out, _ = execNodeCLI("--json", "node", "mempool")
	if !strings.Contains(out, "1") {
		t.Fatalf("mempool output: %s", out)
	}

	out, err = execNodeCLI("--json", "node", "mine")
	if err != nil {
		t.Fatalf("mine failed: %v", err)
	}
	if !strings.Contains(out, "hash") {
		t.Fatalf("unexpected mine output: %s", out)
	}

	out, _ = execNodeCLI("--json", "node", "mempool")
	if !strings.Contains(out, "0") {
		t.Fatalf("mempool not emptied: %s", out)
	}
}

// execNodeCLI executes the root command with args capturing stdout.
func execNodeCLI(args ...string) (string, error) {
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
