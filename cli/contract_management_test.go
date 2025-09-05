package cli

import (
	"strings"
	"testing"
)

// TestContractManagerInfoError checks that querying a missing contract prints an error.
func TestContractManagerInfoError(t *testing.T) {
	out, err := execCommand("contract-mgr", "info", "missing")
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
	if !strings.Contains(out, "error") {
		t.Fatalf("expected error output, got %q", out)
	}
}
