package cli

import (
	"strings"
	"testing"
)

func TestFeesEstimate(t *testing.T) {
	out, err := execCommand("fees", "estimate", "--type", "transfer", "--units", "1")
	if err != nil {
		t.Fatalf("estimate failed: %v", err)
	}
	if !strings.Contains(out, "Total") {
		t.Fatalf("unexpected output: %s", out)
	}
}
