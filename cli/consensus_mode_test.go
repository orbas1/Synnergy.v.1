package cli

import (
	"strings"
	"testing"
)

// TestConsensusModeSetShow verifies set and show commands.
func TestConsensusModeSetShow(t *testing.T) {
	if _, err := execCommand("consensus-mode", "set", "pos"); err != nil {
		t.Fatalf("set failed: %v", err)
	}
	out, err := execCommand("consensus-mode", "show")
	if err != nil {
		t.Fatalf("show failed: %v", err)
	}
	if !strings.Contains(strings.ToLower(out), "pos") {
		t.Fatalf("unexpected output: %s", out)
	}
}
