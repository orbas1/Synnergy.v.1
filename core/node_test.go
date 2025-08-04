package core

import "testing"

func TestMineBlockFeeDistribution(t *testing.T) {
	ledger := NewLedger()
	ledger.Credit("alice", 200)
	node := NewNode("miner", "addr", ledger)
	validator := "validator1"
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
	expected := ShareProportional(AdjustForBlockUtilization(dist.ValidatorsMiners, 1, node.MaxTxPerBlock), map[string]uint64{validator: 2, "miner": 1})
	if ledger.GetBalance(validator) != expected[validator] || ledger.GetBalance("miner") != expected["miner"] {
		t.Fatalf("unexpected shares: got validator %d miner %d", ledger.GetBalance(validator), ledger.GetBalance("miner"))
	}
}
