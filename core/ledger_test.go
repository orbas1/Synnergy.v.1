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

func TestLedgerUTXOAndPool(t *testing.T) {
	l := NewLedger()
	l.Mint("alice", 50)
	if utxos := l.GetUTXOs("alice"); len(utxos) != 1 || utxos[0].Amount != 50 {
		t.Fatalf("unexpected utxo state: %+v", utxos)
	}
	tx := NewTransaction("alice", "bob", 20, 0, 0)
	l.AddToPool(tx)
	if pool := l.Pool(); len(pool) != 1 {
		t.Fatalf("expected pool size 1 got %d", len(pool))
	}
	if err := l.ApplyTransaction(tx); err != nil {
		t.Fatalf("apply failed: %v", err)
	}
	if utxos := l.GetUTXOs("alice"); len(utxos) != 1 || utxos[0].Amount != 30 {
		t.Fatalf("unexpected alice utxo: %+v", utxos)
	}
	if utxos := l.GetUTXOs("bob"); len(utxos) != 1 || utxos[0].Amount != 20 {
		t.Fatalf("unexpected bob utxo: %+v", utxos)
	}
}

func TestLedgerValidation(t *testing.T) {
	l := NewLedger()
	if err := l.AddBlock(nil); err != ErrNilBlock {
		t.Fatalf("expected ErrNilBlock got %v", err)
	}
	if err := l.Transfer("", "bob", 10, 0); err != ErrEmptyAddress {
		t.Fatalf("expected ErrEmptyAddress got %v", err)
	}
	if err := l.ApplyTransaction(nil); err != ErrNilTransaction {
		t.Fatalf("expected ErrNilTransaction got %v", err)
	}
}

func TestLedgerBurn(t *testing.T) {
	l := NewLedger()
	l.Credit("alice", 500)

	if err := l.Burn("", 10); err != ErrEmptyAddress {
		t.Fatalf("expected ErrEmptyAddress got %v", err)
	}

	if err := l.Burn("alice", 0); err == nil {
		t.Fatalf("expected validation error for zero amount")
	}

	if err := l.Burn("alice", 600); err == nil {
		t.Fatalf("expected insufficient funds error")
	}

	if err := l.Burn("alice", 200); err != nil {
		t.Fatalf("burn failed: %v", err)
	}

	if bal := l.GetBalance("alice"); bal != 300 {
		t.Fatalf("unexpected balance after burn: %d", bal)
	}
}
