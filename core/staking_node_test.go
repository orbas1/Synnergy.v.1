package core

import "testing"

func TestStakingNodeStakeAndUnstake(t *testing.T) {
	sn := NewStakingNode()
	sn.Stake("addr1", 100)
	if sn.Balance("addr1") != 100 {
		t.Fatalf("expected 100, got %d", sn.Balance("addr1"))
	}
	sn.Unstake("addr1", 40)
	if sn.Balance("addr1") != 60 {
		t.Fatalf("expected 60 after unstake, got %d", sn.Balance("addr1"))
	}
	sn.Unstake("addr1", 100)
	if sn.Balance("addr1") != 0 {
		t.Fatalf("expected 0 after removing stake, got %d", sn.Balance("addr1"))
	}
}

func TestStakingNodeTotal(t *testing.T) {
	sn := NewStakingNode()
	sn.Stake("a", 50)
	sn.Stake("b", 75)
	if sn.TotalStaked() != 125 {
		t.Fatalf("expected total 125, got %d", sn.TotalStaked())
	}
}
