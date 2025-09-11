package cli

import (
	"strings"
	"testing"
)

// TestPeerConnectGas ensures peer connection emits a gas cost line.
func TestPeerConnectGas(t *testing.T) {
	out, err := execCommand("peer", "connect", "addr1")
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") {
		t.Fatalf("expected gas cost, got %q", out)
	}
}

// TestPeerCountGas ensures counting peers emits a gas cost line and output.
func TestPeerCountGas(t *testing.T) {
	out, err := execCommand("peer", "count")
	if err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "count") {
		t.Fatalf("expected count with gas cost, got %q", out)
	}
}
