package whitepaperdetailed

import (
        "math"
        "testing"
)

func TestSynthronCoinSpecification(t *testing.T) {
        spec := SynthronCoin()
        if err := spec.Validate(); err != nil {
                t.Fatalf("expected specification to be valid: %v", err)
        }

        if spec.Symbol != "SYN" {
                t.Fatalf("unexpected symbol: %s", spec.Symbol)
        }
        if spec.MaxSupply != 10_000_000_000 {
                t.Fatalf("unexpected max supply: %d", spec.MaxSupply)
        }
        if spec.InitialSupply >= spec.MaxSupply {
                t.Fatalf("initial supply should be lower than max supply")
        }
        if spec.FeeModel.Allocation.Treasury <= 0 || spec.FeeModel.Allocation.Validators <= 0 {
                t.Fatalf("fee allocations must be positive")
        }
        sum := spec.FeeModel.Allocation.Treasury + spec.FeeModel.Allocation.Validators + spec.FeeModel.Allocation.Burn + spec.FeeModel.Allocation.Ecosystem
        if math.Abs(sum-1.0) > feeTolerance {
                t.Fatalf("fee allocation should sum to 1, got %.6f", sum)
        }

        if len(spec.Emissions) < 2 {
                t.Fatalf("expected multiple emission epochs")
        }
        // verify emissions are sorted
        for i := 1; i < len(spec.Emissions); i++ {
                if spec.Emissions[i-1].Height >= spec.Emissions[i].Height {
                        t.Fatalf("emission epochs not strictly increasing")
                }
        }

        height := spec.Emissions[1].Height
        circulating := spec.CirculatingSupply(height)
        if circulating <= spec.InitialSupply {
                t.Fatalf("circulating supply should increase once emissions unlock")
        }
        remaining := spec.RemainingEmissions(height)
        if remaining == 0 {
                t.Fatalf("remaining emissions should still exist before the final epoch")
        }
        if final := spec.RemainingEmissions(spec.Emissions[len(spec.Emissions)-1].Height); final != 0 {
                t.Fatalf("remaining emissions should be zero after the last scheduled epoch")
        }

        next, ok := spec.NextEmissionAfter(spec.Emissions[0].Height)
        if !ok || next.Height != spec.Emissions[1].Height {
                t.Fatalf("expected to find the next emission epoch")
        }

        snapshot := spec.Snapshot(height)
        if snapshot.CirculatingSupply != circulating {
                t.Fatalf("snapshot supply mismatch")
        }
        ratio := spec.EmissionRatio(height)
        if ratio <= 0 || ratio >= 1 {
                t.Fatalf("emission ratio should be between 0 and 1, got %.4f", ratio)
        }
}
