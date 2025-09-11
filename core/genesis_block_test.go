package core

import "testing"

func TestInitGenesis(t *testing.T) {
	ledger := NewLedger()
	node := NewNode("n1", "addr", ledger)
	wallets := DefaultGenesisWallets()
	stats, block, err := node.InitGenesis(wallets)
	if err != nil {
		t.Fatalf("init genesis: %v", err)
	}
	if block == nil || block.Hash == "" {
		t.Fatalf("expected block hash, got %v", block)
	}
	if h, _ := ledger.Head(); h != 1 {
		t.Fatalf("expected ledger height 1, got %d", h)
	}
	if bal := ledger.GetBalance(wallets.CreatorWallet); bal != GenesisAllocation {
		t.Fatalf("creator balance %d", bal)
	}
	if stats.Circulating != GenesisAllocation {
		t.Fatalf("circulating %d", stats.Circulating)
	}
	if stats.Remaining != MaxSupply-GenesisAllocation {
		t.Fatalf("remaining %d", stats.Remaining)
	}
	if _, _, err := node.InitGenesis(wallets); err == nil {
		t.Fatalf("expected error on second init")
	}

	bad := wallets
	bad.CreatorWallet = ""
	if _, _, err := node.InitGenesis(bad); err == nil {
		t.Fatalf("expected error for invalid wallets")
	}
}
