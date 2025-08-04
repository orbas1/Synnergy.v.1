package core

import "testing"

func TestProtocolRegistry(t *testing.T) {
	reg := NewProtocolRegistry()
	p := reg.RegisterProtocol("IBC")
	if _, ok := reg.GetProtocol(p.ID); !ok {
		t.Fatalf("protocol not found")
	}
	if len(reg.ListProtocols()) != 1 {
		t.Fatalf("list: expected 1 protocol")
	}
}
