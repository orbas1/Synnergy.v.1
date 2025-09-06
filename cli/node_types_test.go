package cli

import (
	"strings"
	"testing"
)

// TestNodeAddrCLI validates parse and validate commands produce JSON.
func TestNodeAddrCLI(t *testing.T) {
	out, err := execCommand("nodeaddr", "parse", "n1", "--json")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if !strings.Contains(out, "parsed") {
		t.Fatalf("unexpected parse output: %s", out)
	}

	out, err = execCommand("nodeaddr", "validate", "n1", "--json")
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if !strings.Contains(out, "true") {
		t.Fatalf("unexpected validate output: %s", out)
	}
}
