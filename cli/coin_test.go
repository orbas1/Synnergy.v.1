package cli

import (
	"strings"
	"testing"
)

// TestCoinInfoJSON ensures info subcommand emits structured JSON.
func TestCoinInfoJSON(t *testing.T) {
	out, err := execCommand("coin", "--json", "info")
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
	if !strings.Contains(out, "\"name\"") {
		t.Fatalf("expected JSON output, got %s", out)
	}
}

// TestCoinRewardValidation ensures invalid heights are rejected.
func TestCoinRewardValidation(t *testing.T) {
	if _, err := execCommand("coin", "reward", "abc"); err == nil {
		t.Fatalf("expected error for invalid height")
	}
}
