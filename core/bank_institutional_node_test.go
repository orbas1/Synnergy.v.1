package core

import "testing"

// TestNewBankInstitutionalNode verifies construction of a banking institutional node
// and ensures its embedded Node and state are initialized properly.
func TestNewBankInstitutionalNode(t *testing.T) {
	ledger := NewLedger()
	bn := NewBankInstitutionalNode("bank1", "addr1", ledger)
	if bn.Node == nil {
		t.Fatalf("expected embedded Node to be initialized")
	}
	if bn.ID != "bank1" {
		t.Fatalf("expected ID 'bank1', got %s", bn.ID)
	}
	if bn.Addr != "addr1" {
		t.Fatalf("expected Addr 'addr1', got %s", bn.Addr)
	}
	if bn.Ledger != ledger {
		t.Fatalf("ledger not assigned")
	}
	if bn.Institutions == nil || len(bn.Institutions) != 0 {
		t.Fatalf("institutions map not initialized correctly")
	}
}

// TestBankInstitutionalNodeRegistration exercises institution registration
// behaviour including duplicate registrations and lookup correctness.
func TestBankInstitutionalNodeRegistration(t *testing.T) {
	ledger := NewLedger()
	bn := NewBankInstitutionalNode("bank1", "addr1", ledger)
	if bn.IsRegistered("BankA") {
		t.Fatalf("expected BankA to be unregistered initially")
	}
	bn.RegisterInstitution("BankA")
	bn.RegisterInstitution("BankB")
	bn.RegisterInstitution("BankC")
	if !bn.IsRegistered("BankA") || !bn.IsRegistered("BankB") || !bn.IsRegistered("BankC") {
		t.Fatalf("registration failed for some institutions")
	}
	if len(bn.Institutions) != 3 {
		t.Fatalf("expected 3 registered institutions, got %d", len(bn.Institutions))
	}
	// Registering the same institution again should not change the count.
	bn.RegisterInstitution("BankA")
	if len(bn.Institutions) != 3 {
		t.Fatalf("duplicate registration altered institution map size")
	}
}

// TestBankInstitutionalNodeIsolation ensures institution registrations are
// isolated to the node on which they were performed.
func TestBankInstitutionalNodeIsolation(t *testing.T) {
	ledger := NewLedger()
	bn1 := NewBankInstitutionalNode("bank1", "addr1", ledger)
	bn2 := NewBankInstitutionalNode("bank2", "addr2", ledger)

	bn1.RegisterInstitution("BankA")
	if !bn1.IsRegistered("BankA") {
		t.Fatalf("bn1 expected to register BankA")
	}
	if bn2.IsRegistered("BankA") {
		t.Fatalf("bn2 should not have BankA registered")
	}
}

// TestRegisterEmptyInstitution ensures even empty institution names are handled.
func TestRegisterEmptyInstitution(t *testing.T) {
	ledger := NewLedger()
	bn := NewBankInstitutionalNode("bank1", "addr1", ledger)
	bn.RegisterInstitution("")
	if !bn.IsRegistered("") {
		t.Fatalf("empty institution name should be registered")
	}
	if len(bn.Institutions) != 1 {
		t.Fatalf("expected 1 institution after empty name registration, got %d", len(bn.Institutions))
	}
}
