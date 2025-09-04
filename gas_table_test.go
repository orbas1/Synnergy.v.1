package synnergy

import "testing"

func TestGasTableIncludesNewOpcodes(t *testing.T) {
	ResetGasTable()
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
}
