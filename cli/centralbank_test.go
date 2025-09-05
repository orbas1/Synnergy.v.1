package cli

import (
	"strings"
	"testing"
)

// TestCentralbankInfo ensures the info subcommand emits basic node details.
func TestCentralbankInfo(t *testing.T) {
	out, err := execCommand("centralbank", "info")
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
	if out == "" || out == "\n" {
		t.Fatalf("expected output, got %q", out)
	}
	if !strings.Contains(out, "id:") {
		t.Fatalf("missing id in output: %s", out)
	}
}
