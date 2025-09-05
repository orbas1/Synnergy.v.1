package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

// TestAddressUtilities exercises parse, bytes and short subcommands.
func TestAddressUtilities(t *testing.T) {
	hex := core.AddressZero.Hex()
	if out, err := execCommand("address", "parse", hex); err != nil || out != hex {
		t.Fatalf("parse failed: %v %q", err, out)
	}
	if out, err := execCommand("address", "bytes", hex); err != nil || strings.ToLower(out) != strings.TrimPrefix(hex, "0x") {
		t.Fatalf("bytes failed: %v %q", err, out)
	}
	if out, err := execCommand("address", "short", hex); err != nil || out == "" {
		t.Fatalf("short failed: %v %q", err, out)
	}
}
