package core

import "testing"

func TestDefaultGasTable(t *testing.T) {
	initGasTable()
	if GasCost(OpAdd) != DefaultGasCost {
		t.Fatalf("expected default gas cost %d for OpAdd, got %d", DefaultGasCost, GasCost(OpAdd))
	}
	SetGasCost(OpAdd, 5)
	if GasCost(OpAdd) != 5 {
		t.Fatalf("expected updated gas cost 5 for OpAdd, got %d", GasCost(OpAdd))
	}
}
