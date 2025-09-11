package core

import "testing"

func TestSidechainOps(t *testing.T) {
	reg := NewSidechainRegistry()
	reg.Register("chain1", "meta", nil)
	ops := NewSidechainOps(reg)
	if err := ops.Deposit("chain1", "alice", 100); err != nil {
		t.Fatalf("deposit: %v", err)
	}
	bal, err := ops.EscrowBalance("chain1", "alice")
	if err != nil || bal != 100 {
		t.Fatalf("balance: %d %v", bal, err)
	}
	if err := ops.Withdraw("chain1", "alice", 60, "proof"); err != nil {
		t.Fatalf("withdraw: %v", err)
	}
	bal, _ = ops.EscrowBalance("chain1", "alice")
	if bal != 40 {
		t.Fatalf("unexpected balance %d", bal)
	}
	deps, err := ops.ListDeposits("chain1")
	if err != nil || deps["alice"] != 40 {
		t.Fatalf("list deposits: %v", err)
	}
}
