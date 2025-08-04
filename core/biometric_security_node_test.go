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

	// Test SecureExecute with correct and incorrect biometrics
	if err := bsn.SecureExecute(admin, bio, func() error {
		bsn.Node.ID = "updated"
		return nil
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bsn.Node.ID != "updated" {
		t.Fatal("secure execute did not run")
	}
	if err := bsn.SecureExecute(admin, []byte("bad"), nil); err == nil {
		t.Fatal("expected verification failure")
	}

	// Test removal of biometrics
	bsn.Remove(admin)
	if bsn.Auth.Enrolled(admin) {
		t.Fatal("expected biometric data to be removed")
	}
}
