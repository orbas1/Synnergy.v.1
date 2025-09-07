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

// TestMainRunsFirewallCommand ensures the standalone binary delegates to the
// firewall subcommand of the shared CLI.
func TestMainRunsFirewallCommand(t *testing.T) {
	oldArgs := os.Args
	os.Args = []string{"firewall", "check", "1.2.3.4"}
	defer func() { os.Args = oldArgs }()

	out := captureOutput(main)
	if !strings.Contains(out, "true") {
		t.Fatalf("unexpected output: %s", out)
	}
}
