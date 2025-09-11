package core

import "testing"

func TestImmutabilityEnforcer(t *testing.T) {
	genesis := &Block{Hash: "gen"}
	l := NewLedger()
	if err := l.AddBlock(genesis); err != nil {
		t.Fatalf("add block: %v", err)
	}

	enforcer := NewImmutabilityEnforcer(genesis)
	if err := enforcer.CheckLedger(l); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := enforcer.CheckLedger(nil); err != ErrNilLedger {
		t.Fatalf("expected ErrNilLedger, got %v", err)
	}
	empty := NewLedger()
	if err := enforcer.CheckLedger(empty); err != ErrGenesisMissing {
		t.Fatalf("expected ErrGenesisMissing, got %v", err)
	}
	// modify genesis hash to trigger failure
	l.blocks[0].Hash = "other"
	if err := enforcer.CheckLedger(l); err != ErrGenesisChanged {
		t.Fatalf("expected mismatch error")
	}
}
