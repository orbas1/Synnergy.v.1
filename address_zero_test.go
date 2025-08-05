package synnergy

import "testing"

func TestIsZeroAddress(t *testing.T) {
	if !IsZeroAddress(AddressZero) {
		t.Fatalf("expected zero address")
	}
	if IsZeroAddress("0x1") {
		t.Fatalf("non-zero address detected as zero")
	}
}
