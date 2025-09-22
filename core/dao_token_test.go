package core

import "testing"

func TestDAOTokenLedger(t *testing.T) {
	mgr := NewDAOManager()
	mgr.AuthorizeRelayer("admin")
	dao, err := mgr.Create("d1", "admin")
	if err != nil {
		t.Fatalf("create dao: %v", err)
	}
	mgr.AuthorizeRelayer("alice")
	mgr.AuthorizeRelayer("bob")
	if err := mgr.Join(dao.ID, "alice"); err != nil {
		t.Fatalf("join alice: %v", err)
	}
	if err := mgr.Join(dao.ID, "bob"); err != nil {
		t.Fatalf("join bob: %v", err)
	}

	ledger := NewLedger()
	ledger.Mint("admin", 100)
	l := NewDAOTokenLedger(mgr, ledger)

	if err := l.Mint(dao.ID, "alice", "alice", 10); err == nil {
		t.Fatalf("expected unauthorized mint")
	}
	if err := l.Mint(dao.ID, "admin", "carol", 10); err == nil {
		t.Fatalf("expected member missing")
	}
	if err := l.Mint(dao.ID, "admin", "alice", 10); err != nil {
		t.Fatalf("mint: %v", err)
	}
	if bal := l.Balance(dao.ID, "alice"); bal != 10 {
		t.Fatalf("expected 10, got %d", bal)
	}
	if err := l.Transfer(dao.ID, "alice", "bob", 5); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if l.Balance(dao.ID, "alice") != 5 || l.Balance(dao.ID, "bob") != 5 {
		t.Fatalf("unexpected balances")
	}
	if err := l.Transfer(dao.ID, "alice", "carol", 1); err == nil {
		t.Fatalf("expected transfer error")
	}
	if err := l.Burn(dao.ID, "admin", "bob", 3); err != nil {
		t.Fatalf("burn: %v", err)
	}
	if l.Balance(dao.ID, "bob") != 2 {
		t.Fatalf("unexpected burn balance")
	}
	if err := l.Burn(dao.ID, "alice", "bob", 1); err == nil {
		t.Fatalf("expected unauthorized burn")
	}
	if err := l.Burn(dao.ID, "admin", "bob", 5); err == nil {
		t.Fatalf("expected burn error")
	}
}
