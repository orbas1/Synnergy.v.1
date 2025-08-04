package core

import "testing"

// TestValidatorManagerLifecycle verifies basic validator operations.

func TestValidatorManagerLifecycle(t *testing.T) {
	vm := NewValidatorManager(10)
	if err := vm.Add("v1", 5); err == nil {
		t.Fatalf("expected error for stake below minimum")
	}
	if err := vm.Add("v1", 20); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if vm.Stake("v1") != 20 {
		t.Fatalf("unexpected stake")
	}
	elig := vm.Eligible()
	if _, ok := elig["v1"]; !ok {
		t.Fatalf("validator not eligible")
	}
	vm.Slash("v1")
	if vm.Stake("v1") != 10 {
		t.Fatalf("stake not halved after slash")
	}
	if _, ok := vm.Eligible()["v1"]; ok {
		t.Fatalf("slashed validator should not be eligible")
	}
	vm.Remove("v1")
	if vm.Stake("v1") != 0 {
		t.Fatalf("validator not removed")
	}
}
