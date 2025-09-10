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
	if len(cn.Relayers) != 0 {
		t.Fatalf("expected empty relayer set")
	}
}

// TestCustody verifies that assets are recorded for users.
func TestCustodialNodeCustody(t *testing.T) {
	ledger := NewLedger()
	cn := NewCustodialNode("custodian", "addr", ledger)
	cn.Custody("alice", 100)
	cn.Custody("alice", 50)
	if cn.Balance("alice") != 150 {
		t.Fatalf("unexpected holdings: %v", cn.Balance("alice"))
	}
}

// TestCustodialNodeRelease exercises successful and failing release paths.
func TestCustodialNodeRelease(t *testing.T) {
	ledger := NewLedger()
	cn := NewCustodialNode("custodian", "addr", ledger)
	cn.Custody("alice", 150)

	if err := cn.Release("alice", 40, "relay1"); err == nil {
		t.Fatalf("expected unauthorized release to fail")
	}
	cn.AuthorizeRelayer("relay1")
	if err := cn.Release("alice", 40, "relay1"); err != nil {
		t.Fatalf("expected successful release: %v", err)
	}
	if cn.Balance("alice") != 110 {
		t.Fatalf("unexpected holdings after release: %d", cn.Balance("alice"))
	}
	if bal := ledger.GetBalance("alice"); bal != 40 {
		t.Fatalf("ledger not credited: %d", bal)
	}

	if err := cn.Release("alice", 200, "relay1"); err == nil {
		t.Fatalf("expected release to fail when holdings insufficient")
	}
	if cn.Balance("alice") != 110 {
		t.Fatalf("holdings changed after failed release: %d", cn.Balance("alice"))
	}
	if bal := ledger.GetBalance("alice"); bal != 40 {
		t.Fatalf("ledger balance changed after failed release: %d", bal)
	}

	if err := cn.Release("alice", 110, "relay1"); err != nil {
		t.Fatalf("expected second release to succeed: %v", err)
	}
	if cn.Balance("alice") != 0 {
		t.Fatalf("holdings should be zero: %d", cn.Balance("alice"))
	}
	if bal := ledger.GetBalance("alice"); bal != 150 {
		t.Fatalf("ledger credit incorrect: %d", bal)
	}

	if err := cn.Release("bob", 10, "relay1"); err == nil {
		t.Fatalf("expected release to fail for unknown user")
	}
}
