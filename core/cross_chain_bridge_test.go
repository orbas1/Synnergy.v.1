package core

import "testing"

func TestBridgeTransferManager(t *testing.T) {
	mgr := NewBridgeTransferManager()
	tr, err := mgr.Deposit("bridge1", "alice", "bob", 100, "token")
	if err != nil {
		t.Fatalf("deposit: %v", err)
	}
	if err := mgr.Claim(tr.ID, "proof"); err != nil {
		t.Fatalf("claim: %v", err)
	}
	if tr.Status != "released" {
		t.Fatalf("expected released status")
	}
	if len(mgr.ListTransfers()) != 1 {
		t.Fatalf("list: expected 1 transfer")
	}
}
