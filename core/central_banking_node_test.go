package core

import "testing"

func TestCentralBankingNode(t *testing.T) {
	ledger := NewLedger()
	cbn := NewCentralBankingNode("id", "addr", ledger, "initial")
	if cbn.Node == nil || cbn.MonetaryPolicy != "initial" {
		t.Fatalf("node not initialized correctly")
	}

	cbn.UpdatePolicy("updated")
	if cbn.MonetaryPolicy != "updated" {
		t.Fatalf("policy not updated")
	}

	cbn.Mint("treasury", 100)
	if bal := ledger.GetBalance("treasury"); bal != 100 {
		t.Fatalf("minting failed, balance %d", bal)
	}
}
