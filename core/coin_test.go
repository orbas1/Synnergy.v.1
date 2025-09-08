package core

import (
	"math"
	"testing"
)

func TestBlockRewardHalving(t *testing.T) {
	if BlockReward(0) != InitialBlockReward {
		t.Fatalf("unexpected reward at height 0")
	}
	if BlockReward(HalvingInterval) != InitialBlockReward/2 {
		t.Fatalf("reward not halved at interval")
	}
	if BlockReward(HalvingInterval*63) != 0 {
		t.Fatalf("reward should be zero after 63 halvings")
	}
}

func TestCirculatingAndRemainingSupply(t *testing.T) {
	height := uint64(2)
	expected := GenesisAllocation + BlockReward(0) + BlockReward(1)
	if got := CirculatingSupply(height); got != expected {
		t.Fatalf("circulating supply %d != %d", got, expected)
	}
	if rem := RemainingSupply(height); rem != MaxSupply-expected {
		t.Fatalf("remaining supply %d != %d", rem, MaxSupply-expected)
	}
}

func TestCirculatingSupplyCapped(t *testing.T) {
	if supply := CirculatingSupply(^uint64(0)); supply != MaxSupply {
		t.Fatalf("circulating supply should cap at max, got %d", supply)
	}
}

func TestEconomicHelpers(t *testing.T) {
	if v := InitialPrice(1, 2, 3, 4, 5, 6); math.Abs(v-(1.0+2.0+(3.0*4.0)/5.0)*6.0) > 1e-9 {
		t.Fatalf("initial price mismatch")
	}
	if v := AlphaFactor(1, 2, 3, 4); math.Abs(v-(3.0*1.0+2.0+3.0)/4.0) > 1e-9 {
		t.Fatalf("alpha factor mismatch")
	}
	if v := MinimumStake(100, 10, 5, 2); math.Abs(v-(100.0/(10.0*5.0))*2.0) > 1e-9 {
		t.Fatalf("minimum stake mismatch")
	}
	if MinimumStake(100, 0, 5, 2) != 0 {
		t.Fatalf("minimum stake should be zero when reward or supply zero")
	}
	if v := LockupDuration(10, 2, 4, 0.5); math.Abs(v-(10.0*(2.0/4.0*10.0)+0.5*20.0)) > 1e-9 {
		t.Fatalf("lockup duration mismatch")
	}
	if LockupDuration(10, 2, 0, 0.5) != 10 {
		t.Fatalf("lockup duration should return base when threshold zero")
	}
	ratio := PriceToSupplyRatio(100, 0)
	if math.Abs(ratio-100.0/float64(CirculatingSupply(0))) > 1e-9 {
		t.Fatalf("price to supply ratio mismatch")
	}
}
