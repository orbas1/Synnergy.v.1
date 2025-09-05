package cli

import (
	"os"
	"strings"
	"testing"
)

// TestCompressionSaveLoad verifies snapshot save/load cycle with JSON output.
func TestCompressionSaveLoad(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "snap-*.bin")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	f.Close()

	out, err := execCommand("compression", "--json", "save", f.Name())
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}
	if !strings.Contains(out, "\"status\":\"saved\"") {
		t.Fatalf("unexpected output: %s", out)
	}

	out, err = execCommand("compression", "--json", "load", f.Name())
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if !strings.Contains(out, "\"height\"") {
		t.Fatalf("unexpected output: %s", out)
	}
}
