package core

import "testing"

func TestIsZeroAddress(t *testing.T) {
	if !IsZeroAddress(AddressZero) {
		t.Fatalf("expected zero address to match")
	}
	if IsZeroAddress("0xabc") {
		t.Fatalf("non-zero address detected as zero")
	}
}
