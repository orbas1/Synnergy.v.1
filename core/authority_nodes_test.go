package core

import "testing"

func TestAuthorityNodeRegistry(t *testing.T) {
	reg := NewAuthorityNodeRegistry()
	if _, err := reg.Register("addr1", "validator"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if !reg.IsAuthorityNode("addr1") {
		t.Fatalf("expected addr1 to be authority node")
	}
	if err := reg.Vote("voter", "addr1"); err != nil {
		t.Fatalf("vote: %v", err)
	}
	elect := reg.Electorate(1)
	if len(elect) != 1 || elect[0] != "addr1" {
		t.Fatalf("unexpected electorate: %v", elect)
	}
	reg.Deregister("addr1")
	if reg.IsAuthorityNode("addr1") {
		t.Fatalf("deregister failed")
	}
}
