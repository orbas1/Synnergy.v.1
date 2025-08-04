package core

import "testing"

func TestCCSRegistry(t *testing.T) {
	reg := NewCCSRegistry()
	n := reg.RegisterNetwork("PoW", "PoS")
	if _, ok := reg.GetNetwork(n.ID); !ok {
		t.Fatalf("network not found")
	}
	if len(reg.ListNetworks()) != 1 {
		t.Fatalf("list: expected 1 network")
	}
}
