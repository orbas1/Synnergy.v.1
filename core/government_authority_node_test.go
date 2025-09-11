package core

import (
	"errors"
	"testing"
)

func TestGovernmentAuthorityNode(t *testing.T) {
	gn := NewGovernmentAuthorityNode("addr", "validator", "dept")
	if gn.Department != "dept" {
		t.Fatalf("unexpected department")
	}
	if gn.Address != "addr" {
		t.Fatalf("unexpected address")
	}
	if err := gn.MintSYN("bob", 1); !errors.Is(err, ErrGovernmentMint) {
		t.Fatalf("unexpected mint error: %v", err)
	}
	if err := gn.UpdateMonetaryPolicy("expansion"); !errors.Is(err, ErrGovernmentPolicy) {
		t.Fatalf("unexpected policy error: %v", err)
	}
}
