package core

import "testing"

func TestLedgerApplyTransaction(t *testing.T) {
	l := NewLedger()
	l.Credit("alice", 100)
	tx := NewTransaction("alice", "bob", 40, 2, 0)
	if err := l.ApplyTransaction(tx); err != nil {
		t.Fatalf("apply failed: %v", err)
	}
	if bal := l.GetBalance("alice"); bal != 58 {
		t.Fatalf("unexpected balance: %d", bal)
	}
	if bal := l.GetBalance("bob"); bal != 40 {
		t.Fatalf("unexpected recipient balance: %d", bal)
	}
	tx2 := NewTransaction("alice", "bob", 100, 1, 1)
	if err := l.ApplyTransaction(tx2); err == nil {
		t.Fatalf("expected insufficient funds error")
	}
}
