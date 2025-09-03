package core

import (
	"context"
	"testing"

	ierr "synnergy/internal/errors"
)

// TestValidatorManagerLifecycle verifies basic validator operations.

func TestValidatorManagerLifecycle(t *testing.T) {
	vm := NewValidatorManager(10)
	if err := vm.Add(context.Background(), "v1", 5); err == nil || !ierr.IsCode(err, ierr.Invalid) {
		t.Fatalf("expected invalid stake error")
	}
	if err := vm.Add(context.Background(), "v1", 20); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if vm.Stake("v1") != 20 {
		t.Fatalf("unexpected stake")
	}
	elig := vm.Eligible()
	if _, ok := elig["v1"]; !ok {
		t.Fatalf("validator not eligible")
	}
	vm.Slash(context.Background(), "v1")
	if vm.Stake("v1") != 10 {
		t.Fatalf("stake not halved after slash")
	}
	if _, ok := vm.Eligible()["v1"]; ok {
		t.Fatalf("slashed validator should not be eligible")
	}
	vm.Remove(context.Background(), "v1")
	if vm.Stake("v1") != 0 {
		t.Fatalf("validator not removed")
	}
}
