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
