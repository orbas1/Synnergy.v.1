package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

// TestSyn5000Index ensures tokens registered via the main command are surfaced
// through the index helpers.
func TestSyn5000Index(t *testing.T) {
	syn5000Tokens = map[string]*core.SYN5000Token{}
	syn5000TokenIndex = core.NewSYN5000Index()

	if _, err := execCommand("syn5000", "create", "--name", "Gamble", "--symbol", "GMB"); err != nil {
		t.Fatalf("create: %v", err)
	}

	if out, err := execCommand("syn5000_index", "list"); err != nil {
		t.Fatalf("list: %v", err)
	} else if !strings.Contains(out, "GMB") {
		t.Fatalf("expected GMB in list, got %s", out)
	}

	if out, err := execCommand("syn5000_index", "summary"); err != nil {
		t.Fatalf("summary: %v", err)
	} else if !strings.Contains(out, "GMB") {
		t.Fatalf("expected GMB in summary, got %s", out)
	}

	if out, err := execCommand("syn5000_index", "detail", "--id", "GMB"); err != nil {
		t.Fatalf("detail: %v", err)
	} else if !strings.Contains(out, "Token Gamble") {
		t.Fatalf("expected detail output, got %s", out)
	}
}
