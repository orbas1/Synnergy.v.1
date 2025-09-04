package core

import "testing"

func TestGovernmentAuthorityNode(t *testing.T) {
	gn := NewGovernmentAuthorityNode("addr", "validator", "dept")
	if gn.Department != "dept" {
		t.Fatalf("unexpected department")
	}
	if gn.Address != "addr" {
		t.Fatalf("unexpected address")
	}
	if err := gn.MintSYN("bob", 1); err == nil {
		t.Fatalf("expected mint restriction")
	}
	if err := gn.UpdateMonetaryPolicy("expansion"); err == nil {
		t.Fatalf("expected policy restriction")
	}
}
