package core

import "testing"

func TestBridgeRegistry(t *testing.T) {
	reg := NewBridgeRegistry()
	b, err := reg.RegisterBridge("chainA", "chainB", "relayer1")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, ok := reg.GetBridge(b.ID); !ok {
		t.Fatalf("bridge not found")
	}
	if err := reg.AuthorizeRelayer(b.ID, "relayer2"); err != nil {
		t.Fatalf("auth: %v", err)
	}
	if err := reg.RevokeRelayer(b.ID, "relayer1"); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if len(reg.ListBridges()) != 1 {
		t.Fatalf("list: expected 1 bridge")
	}
}
