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
	if !reg.IsRelayerAuthorized(b.ID, "relayer2") {
		t.Fatalf("relayer2 should be authorized")
	}
	if reg.IsRelayerAuthorized(b.ID, "relayer1") {
		t.Fatalf("relayer1 should not be authorized")
	}
	if err := reg.RemoveBridge(b.ID); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if len(reg.ListBridges()) != 0 {
		t.Fatalf("list: expected 0 bridges after removal")
	}
}
