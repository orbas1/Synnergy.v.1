package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"synnergy/core"
)

// TestMiningNodeCLI verifies mining start/stop and hashing.
func TestMiningNodeCLI(t *testing.T) {
	miningNode = core.NewMiningNode(1000)

	out, err := execMiningCLI("mining", "start")
	if err != nil || !strings.Contains(out, "started") {
		t.Fatalf("start: %v output %s", err, out)
	}

	out, err = execMiningCLI("--json", "mining", "status")
	if err != nil || !strings.Contains(out, "true") {
		t.Fatalf("status: %v output %s", err, out)
	}

	out, err = execMiningCLI("--json", "mining", "mine", "data")
	if err != nil || !strings.Contains(out, "hash") {
		t.Fatalf("mine: %v output %s", err, out)
	}

	out, err = execMiningCLI("--json", "mining", "mine-until", "data", "0", "--timeout", "1")
	if err != nil || !strings.Contains(out, "hash") {
		t.Fatalf("mine-until: %v output %s", err, out)
	}

	out, err = execMiningCLI("mining", "stop")
	if err != nil || !strings.Contains(out, "stopped") {
		t.Fatalf("stop: %v output %s", err, out)
	}
}

// execMiningCLI executes the CLI and captures stdout.
func execMiningCLI(args ...string) (string, error) {
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
