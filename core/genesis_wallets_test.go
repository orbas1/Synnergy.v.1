package core

import "testing"

func TestDefaultGenesisWalletsDeterministic(t *testing.T) {
	w1 := DefaultGenesisWallets()
	w2 := DefaultGenesisWallets()
	if w1 != w2 {
		t.Fatalf("wallets not deterministic: %#v vs %#v", w1, w2)
	}
	if w1.Genesis == "" || w1.CreatorWallet == "" {
		t.Fatalf("wallet addresses should not be empty: %#v", w1)
	}
}

func TestAllocateToGenesisWallets(t *testing.T) {
	wallets := DefaultGenesisWallets()
	total := uint64(100)
	alloc := AllocateToGenesisWallets(total, wallets)
	dist := DistributeFees(total)

	tests := []struct {
		addr string
		want uint64
	}{
		{wallets.InternalDevelopment, dist.InternalDevelopment},
		{wallets.InternalCharity, dist.InternalCharity},
		{wallets.ExternalCharity, dist.ExternalCharity},
		{wallets.LoanPool, dist.LoanPool},
		{wallets.PassiveIncome, dist.PassiveIncome},
		{wallets.ValidatorsMiners, dist.ValidatorsMiners},
		{wallets.AuthorityNodes, dist.AuthorityNodes},
		{wallets.NodeHosts, dist.NodeHosts},
		{wallets.CreatorWallet, dist.CreatorWallet},
	}

	var sum uint64
	for _, tt := range tests {
		if got := alloc[tt.addr]; got != tt.want {
			t.Fatalf("allocation mismatch for %s: got %d want %d", tt.addr, got, tt.want)
		}
		sum += alloc[tt.addr]
	}
	expected := dist.InternalDevelopment + dist.InternalCharity + dist.ExternalCharity + dist.LoanPool + dist.PassiveIncome + dist.ValidatorsMiners + dist.AuthorityNodes + dist.NodeHosts + dist.CreatorWallet
	if sum != expected {
		t.Fatalf("total allocation mismatch: got %d want %d", sum, expected)
	}
}
