package cli

import (
	"strings"
	"testing"
)

// TestNodeAdapterInfo ensures the adapter reports the node ID with gas.
func TestNodeAdapterInfo(t *testing.T) {
	out, err := execCommand("node_adapter", "info", "--json")
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "id") {
		t.Fatalf("unexpected output: %s", out)
	}
}
