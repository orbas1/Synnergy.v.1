package core

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"synnergy/internal/tokens"
)

func TestBankNodes(t *testing.T) {
	ledger := NewLedger()
	pub, priv, _ := ed25519.GenerateKey(nil)
	addr := hex.EncodeToString(pub)
	b := NewBankInstitutionalNode("id1", addr, ledger)
	if b == nil || b.Node == nil {
		t.Fatal("bank institutional node not created")
	}
	tok := tokens.NewSYN10Token(1, "CBDC", "cSYN", "central", 1, 2)
	c := NewCentralBankingNode("id2", "addr2", ledger, "policy", tok)
	if c.MonetaryPolicy != "policy" {
		t.Fatal("policy not set")
	}
	s := NewCustodialNode("id3", "addr3", ledger)
	s.Custody("user", 10)
	if err := s.Release("user", 5); err != nil {
		t.Fatal("release failed")
	}

	sigReg := ed25519.Sign(priv, []byte("register:inst1"))
	if err := b.RegisterInstitution(addr, "inst1", sigReg, pub); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if !b.IsRegistered("inst1") {
		t.Fatal("institution not registered")
	}
	sigRem := ed25519.Sign(priv, []byte("remove:inst1"))
	if err := b.RemoveInstitution(addr, "inst1", sigRem, pub); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if b.IsRegistered("inst1") {
		t.Fatal("institution still present")
	}
}
