package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// TestMainRunsGovernmentCommand verifies the binary executes the government
// CLI subcommand.
func TestMainRunsGovernmentCommand(t *testing.T) {
	oldArgs := os.Args
	os.Args = []string{"governance", "new", "addr1", "role1", "dept1", "--json"}
	defer func() { os.Args = oldArgs }()

	out := captureOutput(main)
	if !strings.Contains(out, "role1") {
		t.Fatalf("unexpected output: %s", out)
	}
}
