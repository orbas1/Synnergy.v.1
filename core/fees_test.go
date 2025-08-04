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

func TestCalculateBaseFee(t *testing.T) {
	recent := make([]uint64, 1100)
	for i := range recent {
		recent[i] = uint64(i + 1)
	}
	fee := CalculateBaseFee(recent, 0.1)
	if fee != 661 {
		t.Fatalf("unexpected base fee %d", fee)
	}
}

func TestCalculateVariableAndPriorityFee(t *testing.T) {
	if vf := CalculateVariableFee(10, 5); vf != 50 {
		t.Fatalf("variable fee incorrect %d", vf)
	}
	if pf := CalculatePriorityFee(7); pf != 7 {
		t.Fatalf("priority fee incorrect %d", pf)
	}
}

func TestSpecificFeeCalculations(t *testing.T) {
	if fb := FeeForPurchase(3, 1, 2, 4); fb.Total != 11 {
		t.Fatalf("purchase fee incorrect %+v", fb)
	}
	if fb := FeeForTokenUsage(3, 1, 2, 4); fb.Total != 11 {
		t.Fatalf("token fee incorrect %+v", fb)
	}
	if fb := FeeForContract(3, 1, 2, 4); fb.Total != 11 {
		t.Fatalf("contract fee incorrect %+v", fb)
	}
	if fb := FeeForWalletVerification(3, 1, 2, 4); fb.Total != 11 {
		t.Fatalf("verify fee incorrect %+v", fb)
	}
}

func TestFeeForValidatedTransfer(t *testing.T) {
	fb := FeeForValidatedTransfer(10, 1, 2, 3, true)
	if fb.Total != 0 {
		t.Fatalf("validated transfer should be free, got %+v", fb)
	}
	fb = FeeForValidatedTransfer(10, 1, 2, 3, false)
	if fb.Total != 24 {
		t.Fatalf("non-validated transfer fee incorrect %+v", fb)
	}
}
