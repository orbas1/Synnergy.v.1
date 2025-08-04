package tokens

import (
	"encoding/hex"
	"testing"
)

func TestBaseTokenMintTransferBurn(t *testing.T) {
	tok := NewBaseToken(1, "Test", "TST", 0)
	if err := tok.Mint("alice", 100); err != nil {
		t.Fatalf("mint: %v", err)
	}
	if err := tok.Transfer("alice", "bob", 40); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if tok.BalanceOf("bob") != 40 {
		t.Fatalf("expected bob=40 got %d", tok.BalanceOf("bob"))
	}
	if err := tok.Burn("alice", 10); err != nil {
		t.Fatalf("burn: %v", err)
	}
	if tok.TotalSupply() != 90 {
		t.Fatalf("unexpected supply %d", tok.TotalSupply())
	}
}

func TestSYN10Info(t *testing.T) {
	tkn := NewSYN10Token(2, "CBDC", "CBD", "Central Bank", 1.0, 2)
	if err := tkn.Mint("alice", 50); err != nil {
		t.Fatalf("mint: %v", err)
	}
	info := tkn.Info()
	if info.Issuer != "Central Bank" || info.ExchangeRate != 1.0 {
		t.Fatalf("unexpected info %+v", info)
	}
}

func TestSYN1000ReserveValue(t *testing.T) {
	idx := NewSYN1000Index()
	id := idx.Create("Stable", "STBL", 2)
	if err := idx.AddReserve(id, "USD", 100); err != nil {
		t.Fatalf("add reserve: %v", err)
	}
	if err := idx.SetReservePrice(id, "USD", 1.0); err != nil {
		t.Fatalf("set price: %v", err)
	}
	val, err := idx.TotalValue(id)
	if err != nil || val != 100 {
		t.Fatalf("unexpected value %f err %v", val, err)
	}
}

func TestSYN1100Access(t *testing.T) {
	store := NewSYN1100Token()
	data, _ := hex.DecodeString("abcd")
	if err := store.AddRecord(1, "alice", data); err != nil {
		t.Fatalf("add record: %v", err)
	}
	if err := store.GrantAccess(1, "bob"); err != nil {
		t.Fatalf("grant: %v", err)
	}
	if _, err := store.GetRecord(1, "bob"); err != nil {
		t.Fatalf("bob should access: %v", err)
	}
	if err := store.RevokeAccess(1, "bob"); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if _, err := store.GetRecord(1, "bob"); err == nil {
		t.Fatalf("expected access denied")
	}
}
