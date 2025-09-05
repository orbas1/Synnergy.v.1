package cli

import (
	"strings"
	"testing"
)

// TestConnPoolStatsAndClose ensures pool stats are printed and the pool closes.
func TestConnPoolStatsAndClose(t *testing.T) {
	out, err := execCommand("connpool", "stats")
	if err != nil {
		t.Fatalf("stats failed: %v", err)
	}
	if !strings.Contains(out, "active:") {
		t.Fatalf("unexpected output: %s", out)
	}
	if _, err := execCommand("connpool", "close"); err != nil {
		t.Fatalf("close failed: %v", err)
	}
}
