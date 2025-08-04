package core

import "testing"

func TestStakePenaltyManager(t *testing.T) {
	sn := NewStakingNode()
	sn.Stake("validator", 100)
	spm := NewStakePenaltyManager()
	spm.Slash(sn, "validator", 30)
	if sn.Balance("validator") != 70 {
		t.Fatalf("expected 70 after slash, got %d", sn.Balance("validator"))
	}
	spm.Reward(sn, "validator", 10)
	if sn.Balance("validator") != 80 {
		t.Fatalf("expected 80 after reward, got %d", sn.Balance("validator"))
	}
}
