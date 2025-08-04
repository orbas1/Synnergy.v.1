package core

import "testing"

func TestDefaultGasTable(t *testing.T) {
	initGasTable()
	if GasCost(OpAdd) != 0 {
		t.Fatalf("expected zero gas cost for unpriced opcode")
	}
	SetGasCost(OpAdd, 5)
	if GasCost(OpAdd) != 5 {
		t.Fatalf("expected updated gas cost 5 for OpAdd, got %d", GasCost(OpAdd))
	}
}
