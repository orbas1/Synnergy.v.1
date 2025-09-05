package cli

import (
	"strings"
	"testing"
)

// TestConnPoolLifecycle ensures the pool can dial, release and close connections.
func TestConnPoolLifecycle(t *testing.T) {
	if _, err := execCommand("connpool", "dial", "peer1"); err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	out, err := execCommand("connpool", "stats")
	if err != nil {
		t.Fatalf("stats failed: %v", err)
	}
	if !strings.Contains(out, "active: 1") {
		t.Fatalf("expected active 1, got %q", out)
	}
	if _, err := execCommand("connpool", "release", "peer1"); err != nil {
		t.Fatalf("release failed: %v", err)
	}
	out, err = execCommand("connpool", "stats")
	if err != nil {
		t.Fatalf("stats failed: %v", err)
	}
	if !strings.Contains(out, "active: 0") {
		t.Fatalf("expected active 0, got %q", out)
	}
	if _, err := execCommand("connpool", "close"); err != nil {
		t.Fatalf("close failed: %v", err)
	}
}
