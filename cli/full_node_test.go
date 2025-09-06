package cli

import (
	"strings"
	"testing"
)

func TestFullNodeLifecycle(t *testing.T) {
	if _, err := execCommand("fullnode", "create", "--id", "fn1", "--mode", "archive", "--json"); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	out, err := execCommand("fullnode", "set-mode", "pruned", "--json")
	if err != nil {
		t.Fatalf("set-mode failed: %v", err)
	}
	if !strings.Contains(out, "pruned") {
		t.Fatalf("expected pruned mode: %s", out)
	}
}
