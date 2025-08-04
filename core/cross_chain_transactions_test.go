package core

import "testing"

func TestCrossChainTxManager(t *testing.T) {
	l := NewLedger()
	l.Credit("alice", 100)
	l.Credit("bob", 100)
	m := NewCrossChainTxManager(l)
	id1, err := m.LockMint(1, "alice", "charlie", "asset1", 40, "proof")
	if err != nil {
		t.Fatalf("lockmint failed: %v", err)
	}
	if l.GetBalance("alice") != 60 || l.GetBalance("charlie") != 40 {
		t.Fatalf("unexpected balances after lockmint")
	}
	id2, err := m.BurnRelease(1, "bob", "dave", "asset1", 30)
	if err != nil {
		t.Fatalf("burnrelease failed: %v", err)
	}
	if l.GetBalance("bob") != 70 || l.GetBalance("dave") != 30 {
		t.Fatalf("unexpected balances after burnrelease")
	}
	if _, err := m.GetTransfer(id1); err != nil {
		t.Fatalf("get transfer1 failed: %v", err)
	}
	if _, err := m.GetTransfer(id2); err != nil {
		t.Fatalf("get transfer2 failed: %v", err)
	}
	if len(m.ListTransfers()) != 2 {
		t.Fatalf("expected two transfers")
	}
}
