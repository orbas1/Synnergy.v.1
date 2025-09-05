package cli

import (
	"bytes"
	"testing"

	"synnergy/core"
	"synnergy/internal/nodes"
)

// TestBaseNodeLifecycle verifies basic start/stop behaviour via the CLI.
func TestBaseNodeLifecycle(t *testing.T) {
	baseNode = core.NewBaseNode(nodes.Address("base1"))

	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"basenode", "start"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if !baseNode.IsRunning() {
		t.Fatalf("node should be running")
	}

	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"basenode", "stop"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("stop failed: %v", err)
	}
	if baseNode.IsRunning() {
		t.Fatalf("node should be stopped")
	}
}
