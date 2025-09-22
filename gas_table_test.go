package synnergy

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGasTableIncludesNewOpcodes(t *testing.T) {
	ResetGasTable()
	defer ResetGasTable()
	table := LoadGasTable()
	if len(table) == 0 {
		t.Fatalf("expected table entries")
	}
	if !HasOpcode("Security_RaiseAlert") {
		t.Fatalf("missing Security_RaiseAlert opcode")
	}
	if GasCost("Security_RaiseAlert") != 150 {
		t.Fatalf("unexpected cost for Security_RaiseAlert")
	}
	if !HasOpcode("Marketplace_ListContract") {
		t.Fatalf("missing Marketplace_ListContract opcode")
	}
	if GasCost("Marketplace_ListContract") != 80 {
		t.Fatalf("unexpected cost for Marketplace_ListContract")
	}
	if !HasOpcode("RegisterContentNode") {
		t.Fatalf("missing RegisterContentNode opcode")
	}
	if GasCost("RegisterContentNode") != 5 {
		t.Fatalf("unexpected cost for RegisterContentNode")
	}
	if !HasOpcode("MineUntil") {
		t.Fatalf("missing MineUntil opcode")
	}
	if GasCost("MineUntil") != 50 {
		t.Fatalf("unexpected cost for MineUntil")
	}
	if MustGasCost("MineUntil") != 50 {
		t.Fatalf("MustGasCost returned wrong value")
	}
	if !HasOpcode("RegNodeApprove") {
		t.Fatalf("missing RegNodeApprove opcode")
	}
	if GasCost("RegNodeApprove") != 2 {
		t.Fatalf("unexpected cost for RegNodeApprove")
	}
}

func TestRegisterGasCostValidation(t *testing.T) {
	if err := RegisterGasCost("", 1); err == nil {
		t.Fatalf("expected error for empty name")
	}
	if err := RegisterGasCost("Valid", 0); err == nil {
		t.Fatalf("expected error for zero cost")
	}
}

func TestMustGasCostPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for missing opcode")
		}
	}()
	MustGasCost("UnknownOpcode")
}

func TestGasTableSubscription(t *testing.T) {
	ResetGasTable()
	defer ResetGasTable()
	ch, cancel := SubscribeGasTable(1)
	defer cancel()
	snap := <-ch
	if len(snap.Table) == 0 {
		t.Fatalf("expected initial snapshot")
	}
	if err := RegisterGasCost("subscription_test", 77); err != nil {
		t.Fatalf("register cost: %v", err)
	}
	select {
	case snap = <-ch:
		if snap.Table["subscription_test"] != 77 {
			t.Fatalf("expected override propagated, got %v", snap.Table["subscription_test"])
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("timeout waiting for update")
	}
}

func TestGasTableReloadPreservesOverrides(t *testing.T) {
	ResetGasTable()
	defer ResetGasTable()
	if err := RegisterGasCost("override_test", 55); err != nil {
		t.Fatalf("register cost: %v", err)
	}
	snap := ReloadGasTable()
	if snap.Table["override_test"] != 55 {
		t.Fatalf("override not preserved after reload")
	}
}

func TestConcurrentRegisterAndSnapshot(t *testing.T) {
        ResetGasTable()
        defer ResetGasTable()
        var wg sync.WaitGroup
        ctx, cancel := context.WithCancel(context.Background())
        for i := 0; i < 4; i++ {
                wg.Add(1)
                go func(idx int) {
                        defer wg.Done()
                        name := fmt.Sprintf("concurrent_%d", idx)
			for j := 0; j < 10; j++ {
				if err := RegisterGasCost(name, uint64(10+idx+j)); err != nil {
					t.Errorf("register: %v", err)
				}
			}
		}(i)
        }
        wg.Add(1)
        go func() {
                defer wg.Done()
                for {
                        select {
                        case <-ctx.Done():
                                return
                        default:
                                _ = SnapshotGasTable()
                        }
                }
        }()
        time.Sleep(50 * time.Millisecond)
        cancel()
        wg.Wait()
}
