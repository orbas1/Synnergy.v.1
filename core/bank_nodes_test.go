package core

import "testing"

func TestBankNodes(t *testing.T) {
	ledger := NewLedger()
	b := NewBankInstitutionalNode("id1", "addr1", ledger)
	if b == nil || b.Node == nil {
		t.Fatal("bank institutional node not created")
	}
	c := NewCentralBankingNode("id2", "addr2", ledger, "policy")
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
