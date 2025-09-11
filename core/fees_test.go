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
	if dist.InternalDevelopment != 5 || dist.InternalCharity != 5 || dist.ExternalCharity != 5 || dist.LoanPool != 10 || dist.PassiveIncome != 5 || dist.ValidatorsMiners != 59 || dist.AuthorityNodes != 5 || dist.NodeHosts != 5 || dist.CreatorWallet != 1 {
		t.Fatalf("unexpected distribution: %+v", dist)
	}
	total := dist.InternalDevelopment + dist.InternalCharity + dist.ExternalCharity + dist.LoanPool + dist.PassiveIncome + dist.ValidatorsMiners + dist.AuthorityNodes + dist.NodeHosts + dist.CreatorWallet
	if total != 100 {
		t.Fatalf("distribution does not sum to total, got %d", total)
	}
}

func TestDistributeFeesCreatorDisabled(t *testing.T) {
	SetCreatorDistribution(false)
	dist := DistributeFees(100)
	if dist.CreatorWallet != 0 || dist.NodeHosts != 6 {
		t.Fatalf("unexpected distribution when disabled: %+v", dist)
	}
	if !IsCreatorDistributionEnabled() {
		// reset for other tests
		SetCreatorDistribution(true)
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

func TestFeePolicyEnforce(t *testing.T) {
	p := FeePolicy{Cap: 100, Floor: 10}
	if fee, note := p.Enforce(120); fee != 100 || note == "" {
		t.Fatalf("expected cap enforcement, got %d %q", fee, note)
	}
	if fee, note := p.Enforce(5); fee != 10 || note == "" {
		t.Fatalf("expected floor enforcement, got %d %q", fee, note)
	}
	if fee, note := p.Enforce(50); fee != 50 || note != "" {
		t.Fatalf("unexpected change %d %q", fee, note)
	}
}

func TestAdjustFeeRates(t *testing.T) {
	base, variable := AdjustFeeRates(100, 10, 0.5)
	if base != 150 || variable != 15 {
		t.Fatalf("unexpected adjusted rates: base=%d variable=%d", base, variable)
	}
}

func TestEstimateFee(t *testing.T) {
	fb := EstimateFee(TxTypePurchase, 2, 1, 3, 1)
	if fb.Base != 1 || fb.Variable != 6 || fb.Priority != 1 || fb.Total != 8 {
		t.Fatalf("unexpected estimate: %+v", fb)
	}
}

func TestShareProportional(t *testing.T) {
	weights := map[string]uint64{"a": 1, "b": 3}
	shares := ShareProportional(100, weights)
	if shares["a"] != 25 || shares["b"] != 75 {
		t.Fatalf("unexpected shares: %v", shares)
	}
}

func TestAdjustForBlockUtilization(t *testing.T) {
	if v := AdjustForBlockUtilization(100, 95, 100); v != 110 {
		t.Fatalf("high util adjustment failed, got %d", v)
	}
	if v := AdjustForBlockUtilization(100, 40, 100); v != 90 {
		t.Fatalf("low util adjustment failed, got %d", v)
	}
}
