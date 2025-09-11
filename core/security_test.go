package core

import "testing"

func TestSetStakeEnforcesMinimum(t *testing.T) {
	ledger := NewLedger()
	node := NewNode("n1", "addr", ledger)
	if err := node.SetStake("v1", MinStake-1); err == nil {
		t.Fatalf("expected error for stake below minimum")
	}
	if err := node.SetStake("v1", MinStake); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSlashingAndRehabilitation(t *testing.T) {
	ledger := NewLedger()
	node := NewNode("n1", "addr", ledger)
	_ = node.SetStake("v1", MinStake*2)
	node.ReportDoubleSign("v1")
	if _, ok := node.Validators.Eligible()["v1"]; ok {
		t.Fatalf("validator should be slashed")
	}
	if node.Validators.Stake("v1") >= MinStake*2 {
		t.Fatalf("stake not reduced")
	}
	node.Rehabilitate("v1")
	if _, ok := node.Validators.Eligible()["v1"]; !ok {
		t.Fatalf("validator should be rehabilitated")
	}
}

func TestEligibleStakesExcludesSlashed(t *testing.T) {
	ledger := NewLedger()
	node := NewNode("n1", "addr", ledger)
	_ = node.SetStake("v1", MinStake)
	_ = node.SetStake("v2", MinStake)
	node.ReportDowntime("v1")
	elig := node.eligibleStakes()
	if _, ok := elig["v1"]; ok {
		t.Fatalf("slashed validator should not be eligible")
	}
	if _, ok := elig["v2"]; !ok {
		t.Fatalf("validator v2 should be eligible")
	}
}

func TestSubBlockSignature(t *testing.T) {
	tx := NewTransaction("a", "b", 1, 0, 0)
	sb := NewSubBlock([]*Transaction{tx}, "val")
	if !sb.VerifySignature() {
		t.Fatalf("signature verification failed")
	}
}
