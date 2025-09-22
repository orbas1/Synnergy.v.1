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

func TestEnsureGasSchedule(t *testing.T) {
	ResetGasTable()
	schedule := map[string]uint64{
		"EnterpriseBootstrap":     120,
		"EnterpriseConsensusSync": 95,
		"EnterpriseWalletSeal":    60,
	}

	inserted, err := EnsureGasSchedule(schedule)
	if err != nil {
		t.Fatalf("EnsureGasSchedule returned error: %v", err)
	}
	if len(inserted) > len(schedule) {
		t.Fatalf("unexpected inserted values: %v", inserted)
	}

	for name, cost := range schedule {
		if !HasOpcode(name) {
			t.Fatalf("expected opcode %s to be registered", name)
		}
		if GasCost(name) != cost {
			t.Fatalf("unexpected cost for %s: %d", name, GasCost(name))
		}
	}

	// Updating the schedule should refresh costs without duplicating entries.
	updated, err := EnsureGasSchedule(map[string]uint64{"EnterpriseBootstrap": 150})
	if err != nil {
		t.Fatalf("EnsureGasSchedule update returned error: %v", err)
	}
	if len(updated) > 0 {
		t.Fatalf("expected no new insertions, got %v", updated)
	}
	if GasCost("EnterpriseBootstrap") != 150 {
		t.Fatalf("expected updated cost 150, got %d", GasCost("EnterpriseBootstrap"))
	}
}

func TestGasMetadataCatalogue(t *testing.T) {
	ResetGasTable()
	if _, ok := GasMetadataFor("Unknown"); ok {
		t.Fatalf("expected no metadata for unknown opcode")
	}

	if err := RegisterGasMetadata("EnterpriseBootstrap", 200, "orchestrator", "Stage 78 bootstrap sequence"); err != nil {
		t.Fatalf("register metadata: %v", err)
	}
	entry, ok := GasMetadataFor("EnterpriseBootstrap")
	if !ok {
		t.Fatalf("expected metadata for EnterpriseBootstrap")
	}
	if entry.Category != "orchestrator" {
		t.Fatalf("unexpected category: %s", entry.Category)
	}
	if entry.Cost != 200 {
		t.Fatalf("unexpected cost: %d", entry.Cost)
	}

	entries := GasCatalogue()
	if len(entries) == 0 {
		t.Fatalf("expected non-empty metadata catalogue")
	}
	found := false
	for _, e := range entries {
		if e.Name == "EnterpriseBootstrap" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("catalogue missing EnterpriseBootstrap entry")
	}
}
