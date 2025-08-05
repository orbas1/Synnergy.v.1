package core

import "testing"

// TestNewCustodialNode ensures constructor sets up the embedded node and holdings map.
func TestNewCustodialNode(t *testing.T) {
	ledger := NewLedger()
	cn := NewCustodialNode("custodian", "addr", ledger)
	if cn.Node == nil {
		t.Fatalf("embedded node not initialized")
	}
	if cn.Ledger != ledger {
		t.Fatalf("ledger not set")
	}
	if len(cn.Holdings) != 0 {
		t.Fatalf("expected empty holdings")
	}
}

// TestCustody verifies that assets are recorded for users.
func TestCustodialNodeCustody(t *testing.T) {
	ledger := NewLedger()
	cn := NewCustodialNode("custodian", "addr", ledger)
	cn.Custody("alice", 100)
	cn.Custody("alice", 50)
	if cn.Holdings["alice"] != 150 {
		t.Fatalf("unexpected holdings: %v", cn.Holdings["alice"])
	}
}

// TestCustodialNodeRelease exercises successful and failing release paths.
func TestCustodialNodeRelease(t *testing.T) {
	ledger := NewLedger()
	cn := NewCustodialNode("custodian", "addr", ledger)
	cn.Custody("alice", 150)

	if !cn.Release("alice", 40) {
		t.Fatalf("expected successful release")
	}
	if cn.Holdings["alice"] != 110 {
		t.Fatalf("unexpected holdings after release: %d", cn.Holdings["alice"])
	}
	if bal := ledger.GetBalance("alice"); bal != 40 {
		t.Fatalf("ledger not credited: %d", bal)
	}

	if cn.Release("alice", 200) {
		t.Fatalf("expected release to fail when holdings insufficient")
	}
	if cn.Holdings["alice"] != 110 {
		t.Fatalf("holdings changed after failed release: %d", cn.Holdings["alice"])
	}
	if bal := ledger.GetBalance("alice"); bal != 40 {
		t.Fatalf("ledger balance changed after failed release: %d", bal)
	}

	if !cn.Release("alice", 110) {
		t.Fatalf("expected second release to succeed")
	}
	if cn.Holdings["alice"] != 0 {
		t.Fatalf("holdings should be zero: %d", cn.Holdings["alice"])
	}
	if bal := ledger.GetBalance("alice"); bal != 150 {
		t.Fatalf("ledger credit incorrect: %d", bal)
	}

	if cn.Release("bob", 10) {
		t.Fatalf("expected release to fail for unknown user")
	}
}
