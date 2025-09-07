package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDocgenCreatesGuide verifies that running main generates a CLI guide file
// at the location specified by the DOCGEN_OUTPUT environment variable.
func TestDocgenCreatesGuide(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "cli.md")
	t.Setenv("DOCGEN_OUTPUT", out)

	main()

	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(b), "Synnergy blockchain CLI") {
		t.Fatalf("unexpected content: %s", string(b))
	}
}
