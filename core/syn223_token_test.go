package core

import (
	"sync"
	"testing"
)

func TestSYN223TokenTransfer(t *testing.T) {
	token := NewSYN223Token("SYN223", "S223", "alice", 100)
	token.AddToWhitelist("bob")
	if err := token.Transfer("alice", "bob", 40); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if bal := token.BalanceOf("alice"); bal != 60 {
		t.Fatalf("unexpected alice balance %d", bal)
	}
	if bal := token.BalanceOf("bob"); bal != 40 {
		t.Fatalf("unexpected bob balance %d", bal)
	}
}

func TestSYN223TokenWhitelistBlacklist(t *testing.T) {
	token := NewSYN223Token("SYN223", "S223", "owner", 100)
	// Whitelist recipient and perform transfer
	token.AddToWhitelist("carol")
	if err := token.Transfer("owner", "carol", 10); err != nil {
		t.Fatalf("transfer to carol failed: %v", err)
	}
	// Removing from whitelist should block further transfers
	token.RemoveFromWhitelist("carol")
	if err := token.Transfer("owner", "carol", 5); err != ErrRecipientNotWhitelisted {
		t.Fatalf("expected ErrRecipientNotWhitelisted, got %v", err)
	}
	// Blacklist recipient
	token.AddToWhitelist("dave")
	token.AddToBlacklist("dave")
	if err := token.Transfer("owner", "dave", 5); err != ErrAddressBlacklisted {
		t.Fatalf("expected ErrAddressBlacklisted, got %v", err)
	}
	// Remove from blacklist to allow transfer
	token.RemoveFromBlacklist("dave")
	if err := token.Transfer("owner", "dave", 5); err != nil {
		t.Fatalf("transfer to dave after unblacklist failed: %v", err)
	}
	// Blacklist sender should block transfer
	token.AddToWhitelist("eve")
	token.AddToBlacklist("owner")
	if err := token.Transfer("owner", "eve", 5); err != ErrAddressBlacklisted {
		t.Fatalf("expected ErrAddressBlacklisted, got %v", err)
	}
	token.RemoveFromBlacklist("owner")
	if err := token.Transfer("owner", "eve", 5); err != nil {
		t.Fatalf("transfer after removing sender blacklist failed: %v", err)
	}
}

func TestSYN223TokenInsufficientBalance(t *testing.T) {
	token := NewSYN223Token("SYN223", "S223", "alice", 10)
	token.AddToWhitelist("bob")
	if err := token.Transfer("alice", "bob", 20); err != ErrInsufficientBalance {
		t.Fatalf("expected ErrInsufficientBalance, got %v", err)
	}
}

func TestSYN223TokenConcurrentTransfers(t *testing.T) {
	token := NewSYN223Token("SYN223", "S223", "alice", 100)
	token.AddToWhitelist("bob")
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = token.Transfer("alice", "bob", 1)
		}()
	}
	wg.Wait()
	if token.BalanceOf("alice")+token.BalanceOf("bob") != 100 {
		t.Fatalf("balances corrupted")
	}
}
