package cli

import (
	"strings"
	"testing"
)

// TestSimpleVMStatus ensures VM lifecycle commands return expected output.
func TestSimpleVMStatus(t *testing.T) {
	if _, err := execCommand("simplevm", "create", "--json"); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := execCommand("simplevm", "start", "--json"); err != nil {
		t.Fatalf("start: %v", err)
	}
	out, err := execCommand("simplevm", "status", "--json")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "true") {
		t.Fatalf("unexpected output: %s", out)
	}
}
