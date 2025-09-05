package cli

import (
	"testing"

	"synnergy/core"
)

// TestAddressZeroHelpers ensures zero address utilities respond as expected.
func TestAddressZeroHelpers(t *testing.T) {
	if out, err := execCommand("addrzero", "show"); err != nil || out != core.AddressZero.Hex() {
		t.Fatalf("show failed: %v %q", err, out)
	}
	if out, err := execCommand("addrzero", "is", core.AddressZero.Hex()); err != nil || out != "true" {
		t.Fatalf("is failed: %v %q", err, out)
	}
}
