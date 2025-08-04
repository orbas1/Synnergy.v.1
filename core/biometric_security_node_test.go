package core

import "testing"

func TestBiometricSecurityNode(t *testing.T) {
	ledger := NewLedger()
	ledger.Credit("from", 100)
	base := NewNode("node1", "addr1", ledger)
	bsn := NewBiometricSecurityNode(base, nil)

	admin := "admin"
	bio := []byte("admin-bio")
	bsn.Enroll(admin, bio)

	tx := NewTransaction("from", "to", 1, 0, 0)
	if err := bsn.SecureAddTransaction(admin, bio, tx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(base.Mempool) != 1 {
		t.Fatal("transaction not added to mempool")
	}

	tx2 := NewTransaction("from", "to", 1, 0, 1)
	if err := bsn.SecureAddTransaction(admin, []byte("wrong"), tx2); err == nil {
		t.Fatal("expected authentication failure")
	}
	if len(base.Mempool) != 1 {
		t.Fatal("unexpected transaction added")
	}
}
