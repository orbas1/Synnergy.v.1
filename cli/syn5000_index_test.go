package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

// TestSyn5000Index ensures tokens can be added and listed.
func TestSyn5000Index(t *testing.T) {
	gamblingIndex = map[string]core.GamblingToken{}

	if _, err := execCommand("syn5000_index", "add", "GMB"); err != nil {
		t.Fatalf("add: %v", err)
	}
	out, err := execCommand("syn5000_index", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "GMB") {
		t.Fatalf("expected GMB in list, got %s", out)
	}
}
