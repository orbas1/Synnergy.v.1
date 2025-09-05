package cli

import "testing"

// TestBioAuthListEmpty confirms the bioauth list subcommand returns no entries by default.
func TestBioAuthListEmpty(t *testing.T) {
	out, err := execCommand("bioauth", "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if out != "" {
		t.Fatalf("expected no output, got %q", out)
	}
}
