package cli

import (
	"testing"

	"synnergy/core"
)

// TestAccessListEmpty ensures listing roles for an address returns no output by default.
func TestAccessListEmpty(t *testing.T) {
	addr := core.AddressZero.Hex()
	out, err := execCommand("access", "list", addr)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if out != "" {
		t.Fatalf("expected empty list, got %q", out)
	}
}
