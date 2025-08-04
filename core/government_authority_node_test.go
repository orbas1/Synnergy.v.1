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
}
