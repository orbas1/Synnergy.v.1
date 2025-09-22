package core

import (
	"sync"
	"testing"
)

func TestMineBlockFeeDistribution(t *testing.T) {
	ledger := NewLedger()
	ledger.Credit("alice", 200)
	node := NewNode("miner", "addr", ledger)
	validatorWallet, err := NewWallet()
	if err != nil {
		t.Fatalf("wallet: %v", err)
	}
	if err := node.RegisterValidatorWallet(validatorWallet); err != nil {
		t.Fatalf("register validator: %v", err)
	}
	validator := validatorWallet.Address
	node.SetStake(validator, 2)
	tx := NewTransaction("alice", "bob", 10, 100, 0)
	if err := node.AddTransaction(tx); err != nil {
		t.Fatalf("add tx: %v", err)
	}
	block := node.MineBlock()
	if block == nil {
		t.Fatalf("block not mined")
	}
	if bal := ledger.GetBalance("bob"); bal != 10 {
		t.Fatalf("recipient balance %d", bal)
	}
	dist := DistributeFees(100)
	expected := ShareProportional(AdjustForBlockUtilization(dist.ValidatorsMiners, 1, node.MaxTxPerBlock), map[string]uint64{validator: 3, "miner": 1})
	if ledger.GetBalance(validator) != expected[validator] || ledger.GetBalance("miner") != expected["miner"] {
		t.Fatalf("unexpected shares: got validator %d miner %d", ledger.GetBalance(validator), ledger.GetBalance("miner"))
	}
}

// TestNodeConcurrentAddTransaction ensures AddTransaction is safe for concurrent
// use and properly records all transactions.
func TestNodeConcurrentAddTransaction(t *testing.T) {
	ledger := NewLedger()
	ledger.Credit("alice", 1000)
	node := NewNode("n1", "addr", ledger)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = node.AddTransaction(NewTransaction("alice", "b", 1, 1, 0))
		}()
	}
	wg.Wait()
	if len(node.Mempool) != 10 {
		t.Fatalf("expected 10 transactions, got %d", len(node.Mempool))
	}
}
