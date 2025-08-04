package core

import "testing"

func TestBridgeManager(t *testing.T) {
	l := NewLedger()
	l.Credit("alice", 100)
	bm := NewBridgeManager(l)
	bridgeID := bm.RegisterBridge("chainA", "chainB", "relayer1")
	if bridgeID == 0 {
		t.Fatalf("expected bridge id")
	}
	transferID, err := bm.Deposit(bridgeID, "alice", "bob", 50, "token")
	if err != nil {
		t.Fatalf("deposit failed: %v", err)
	}
	if l.GetBalance("alice") != 50 {
		t.Fatalf("expected alice balance 50")
	}
	if err := bm.Claim(transferID, "proof"); err != nil {
		t.Fatalf("claim failed: %v", err)
	}
	if l.GetBalance("bob") != 50 {
		t.Fatalf("expected bob balance 50")
	}
	if _, err := bm.GetTransfer(transferID); err != nil {
		t.Fatalf("get transfer failed: %v", err)
	}
	if len(bm.ListTransfers()) != 1 {
		t.Fatalf("unexpected transfer list length")
	}
}
