package core

import (
	"testing"

	"synnergy/internal/tokens"
)

func TestBankNodes(t *testing.T) {
	ledger := NewLedger()
	b := NewBankInstitutionalNode("id1", "addr1", ledger)
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

	b.RegisterInstitution("inst1")
	if !b.IsRegistered("inst1") {
		t.Fatal("institution not registered")
	}
	b.RemoveInstitution("inst1")
	if b.IsRegistered("inst1") {
		t.Fatal("institution still present")
	}
}
