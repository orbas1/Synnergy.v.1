package cli

import (
	"strings"
	"testing"
)

// TestOpcodesHex ensures the hex subcommand emits gas metrics and output.
func TestOpcodesHex(t *testing.T) {
	out, err := execCommand("opcodes", "hex", "AddBlock")
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") {
		t.Fatalf("expected gas cost, got %q", out)
	}
}
