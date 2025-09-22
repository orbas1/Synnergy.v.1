package tokens

import "testing"

func TestSYN3500TokenLifecycle(t *testing.T) {
	token := NewSYN3500Token("SynStable", "SYNUSD", "Synnergy", 1.0)
	if err := token.Mint("alice", 0); err != ErrInvalidAmount {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
	if err := token.Mint("alice", 100); err != nil {
		t.Fatalf("mint: %v", err)
	}
	if err := token.Transfer("alice", "bob", 50); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if token.BalanceOf("alice") != 50 || token.BalanceOf("bob") != 50 {
		t.Fatalf("unexpected balances: alice %d bob %d", token.BalanceOf("alice"), token.BalanceOf("bob"))
	}
	if err := token.Approve("alice", "bob", 20); err != nil {
		t.Fatalf("approve: %v", err)
	}
	if token.Allowance("alice", "bob") != 20 {
		t.Fatalf("unexpected allowance: %d", token.Allowance("alice", "bob"))
	}
	if err := token.TransferFrom("alice", "bob", "carol", 15); err != nil {
		t.Fatalf("transferFrom: %v", err)
	}
	if token.BalanceOf("carol") != 15 {
		t.Fatalf("unexpected carol balance: %d", token.BalanceOf("carol"))
	}
	token.FreezeAccount("alice")
	if err := token.Transfer("alice", "bob", 1); err != ErrAccountFrozen {
		t.Fatalf("expected ErrAccountFrozen, got %v", err)
	}
	token.UnfreezeAccount("alice")
	if err := token.Redeem("alice", 10); err != nil {
		t.Fatalf("redeem: %v", err)
	}
	if token.TotalSupply() != 90 {
		t.Fatalf("unexpected total supply: %d", token.TotalSupply())
	}
}

func TestSYN3500SetRate(t *testing.T) {
	token := NewSYN3500Token("SynStable", "SYNUSD", "Synnergy", 1.0)
	if err := token.SetRate(0); err != ErrInvalidForexRate {
		t.Fatalf("expected ErrInvalidForexRate, got %v", err)
	}
	if err := token.SetRate(1.05); err != nil {
		t.Fatalf("set rate: %v", err)
	}
	_, _, rate := token.Info()
	if rate != 1.05 {
		t.Fatalf("unexpected rate %f", rate)
	}
}
