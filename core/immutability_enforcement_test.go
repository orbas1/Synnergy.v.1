package core

import "testing"

func TestImmutabilityEnforcer(t *testing.T) {
	genesis := &Block{Hash: "gen"}
	l := NewLedger()
	l.AddBlock(genesis)

	enforcer := NewImmutabilityEnforcer(genesis)
	if err := enforcer.CheckLedger(l); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// modify genesis hash to trigger failure
	l.blocks[0].Hash = "other"
	if err := enforcer.CheckLedger(l); err == nil {
		t.Fatalf("expected mismatch error")
	}
}
