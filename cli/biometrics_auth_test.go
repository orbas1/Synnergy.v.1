package cli

import (
	"strings"
	"testing"
)

// TestBioAuthListEmpty confirms the bioauth list subcommand returns no entries by default.
func TestBioAuthListEmpty(t *testing.T) {
	out, err := execCommand("bioauth", "--json", "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "\"addresses\":[]") {
		t.Fatalf("unexpected output: %s", out)
	}
}
