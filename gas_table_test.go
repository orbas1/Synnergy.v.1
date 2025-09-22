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
	if !HasOpcode("RegisterContentNode") {
		t.Fatalf("missing RegisterContentNode opcode")
	}
	if GasCost("RegisterContentNode") != 5 {
		t.Fatalf("unexpected cost for RegisterContentNode")
	}
	if !HasOpcode("KademliaDistance") {
		t.Fatalf("missing KademliaDistance opcode")
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
	if !HasOpcode("RegNodeFlag") {
		t.Fatalf("missing RegNodeFlag opcode")
	}
	if GasCost("RegNodeFlag") != 1 {
		t.Fatalf("unexpected cost for RegNodeFlag")
	}
	if !HasOpcode("RegNodeLogs") {
		t.Fatalf("missing RegNodeLogs opcode")
	}
	if GasCost("RegNodeLogs") != 1 {
		t.Fatalf("unexpected cost for RegNodeLogs")
	}
	if !HasOpcode("RegNodeAudit") {
		t.Fatalf("missing RegNodeAudit opcode")
	}
	if GasCost("RegNodeAudit") != 3 {
		t.Fatalf("unexpected cost for RegNodeAudit")
	}
	if !HasOpcode("Access_Audit") {
		t.Fatalf("missing Access_Audit opcode")
	}
	if GasCost("Access_Audit") != 2 {
		t.Fatalf("unexpected cost for Access_Audit")
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

func TestEnsureGasCosts(t *testing.T) {
	if err := EnsureGasCosts([]string{"MineBlock", "CreateDAO"}); err != nil {
		t.Fatalf("expected known opcodes: %v", err)
	}
	if err := EnsureGasCosts([]string{"", "DefinitelyMissing"}); err == nil {
		t.Fatalf("expected error for missing opcode")
	}
}
