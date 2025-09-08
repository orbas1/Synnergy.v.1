package core

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"
)

// TestNewBankInstitutionalNode verifies construction of a banking institutional node
// and ensures its embedded Node and state are initialized properly.
func TestNewBankInstitutionalNode(t *testing.T) {
	ledger := NewLedger()
	pub, _, _ := ed25519.GenerateKey(nil)
	bn := NewBankInstitutionalNode("bank1", hex.EncodeToString(pub), ledger)
	if bn.Node == nil {
		t.Fatalf("expected embedded Node to be initialized")
	}
	if bn.ID != "bank1" {
		t.Fatalf("expected ID 'bank1', got %s", bn.ID)
	}
	if bn.Addr != hex.EncodeToString(pub) {
		t.Fatalf("expected Addr %q, got %s", hex.EncodeToString(pub), bn.Addr)
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
	pub, priv, _ := ed25519.GenerateKey(nil)
	addr := hex.EncodeToString(pub)
	bn := NewBankInstitutionalNode("bank1", addr, ledger)
	if bn.IsRegistered("BankA") {
		t.Fatalf("expected BankA to be unregistered initially")
	}
	for _, name := range []string{"BankA", "BankB", "BankC"} {
		sig := ed25519.Sign(priv, []byte("register:"+name))
		if err := bn.RegisterInstitution(addr, name, sig, pub); err != nil {
			t.Fatalf("register %s failed: %v", name, err)
		}
	}
	if !bn.IsRegistered("BankA") || !bn.IsRegistered("BankB") || !bn.IsRegistered("BankC") {
		t.Fatalf("registration failed for some institutions")
	}
	if len(bn.Institutions) != 3 {
		t.Fatalf("expected 3 registered institutions, got %d", len(bn.Institutions))
	}
	// Registering the same institution again should not change the count.
	sig := ed25519.Sign(priv, []byte("register:BankA"))
	if err := bn.RegisterInstitution(addr, "BankA", sig, pub); err != nil {
		t.Fatalf("duplicate registration errored: %v", err)
	}
	if len(bn.Institutions) != 3 {
		t.Fatalf("duplicate registration altered institution map size")
	}
}

// TestBankInstitutionalNodeIsolation ensures institution registrations are
// isolated to the node on which they were performed.
func TestBankInstitutionalNodeIsolation(t *testing.T) {
	ledger := NewLedger()
	pub1, priv1, _ := ed25519.GenerateKey(nil)
	addr1 := hex.EncodeToString(pub1)
	pub2, _, _ := ed25519.GenerateKey(nil)
	addr2 := hex.EncodeToString(pub2)
	bn1 := NewBankInstitutionalNode("bank1", addr1, ledger)
	bn2 := NewBankInstitutionalNode("bank2", addr2, ledger)

	sig := ed25519.Sign(priv1, []byte("register:BankA"))
	if err := bn1.RegisterInstitution(addr1, "BankA", sig, pub1); err != nil {
		t.Fatalf("register failed: %v", err)
	}
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
	pub, priv, _ := ed25519.GenerateKey(nil)
	addr := hex.EncodeToString(pub)
	bn := NewBankInstitutionalNode("bank1", addr, ledger)
	sig := ed25519.Sign(priv, []byte("register:"))
	if err := bn.RegisterInstitution(addr, "", sig, pub); err != nil {
		t.Fatalf("register empty failed: %v", err)
	}
	if !bn.IsRegistered("") {
		t.Fatalf("empty institution name should be registered")
	}
	if len(bn.Institutions) != 1 {
		t.Fatalf("expected 1 institution after empty name registration, got %d", len(bn.Institutions))
	}
}

// TestRemoveInstitution verifies signed removal.
func TestRemoveInstitution(t *testing.T) {
	ledger := NewLedger()
	pub, priv, _ := ed25519.GenerateKey(nil)
	addr := hex.EncodeToString(pub)
	bn := NewBankInstitutionalNode("bank1", addr, ledger)
	sigReg := ed25519.Sign(priv, []byte("register:BankA"))
	if err := bn.RegisterInstitution(addr, "BankA", sigReg, pub); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	sigRem := ed25519.Sign(priv, []byte("remove:BankA"))
	if err := bn.RemoveInstitution(addr, "BankA", sigRem, pub); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if bn.IsRegistered("BankA") {
		t.Fatalf("BankA should be removed")
	}
}
