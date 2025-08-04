package core

import "testing"

func TestCrossChainTxManager(t *testing.T) {
	mgr := NewCrossChainTxManager()
	tx1, err := mgr.LockAndMint("bridge1", "asset1", 50, "proof")
	if err != nil {
		t.Fatalf("lockmint: %v", err)
	}
	tx2, err := mgr.BurnAndRelease("bridge1", "bob", "asset1", 20)
	if err != nil {
		t.Fatalf("burnrelease: %v", err)
	}
	if _, ok := mgr.GetTx(tx1.ID); !ok {
		t.Fatalf("tx1 not found")
	}
	if len(mgr.ListTxs()) != 2 {
		t.Fatalf("list: expected 2 txs")
	}
	if tx2.Type != "burnrelease" {
		t.Fatalf("unexpected type")
	}
}
