package core

import "testing"

func TestFeeForTransfer(t *testing.T) {
	fb := FeeForTransfer(10, 1, 2, 3) // variable=20, total=24
	if fb.Base != 1 || fb.Variable != 20 || fb.Priority != 3 || fb.Total != 24 {
		t.Fatalf("unexpected fee breakdown: %+v", fb)
	}
}

func TestDistributeFees(t *testing.T) {
	dist := DistributeFees(100)
	if dist.InternalDevelopment != 5 || dist.InternalCharity != 5 || dist.ExternalCharity != 5 || dist.LoanPool != 5 || dist.PassiveIncome != 5 || dist.ValidatorsMiners != 69 || dist.NodeHosts != 5 || dist.CreatorWallet != 1 {
		t.Fatalf("unexpected distribution: %+v", dist)
	}
	total := dist.InternalDevelopment + dist.InternalCharity + dist.ExternalCharity + dist.LoanPool + dist.PassiveIncome + dist.ValidatorsMiners + dist.NodeHosts + dist.CreatorWallet
	if total != 100 {
		t.Fatalf("distribution does not sum to total, got %d", total)
	}
}

func TestApplyFeeCapFloor(t *testing.T) {
	if got := ApplyFeeCapFloor(120, 100, 10); got != 100 {
		t.Fatalf("cap not applied, got %d", got)
	}
	if got := ApplyFeeCapFloor(5, 100, 10); got != 10 {
		t.Fatalf("floor not applied, got %d", got)
	}
}
