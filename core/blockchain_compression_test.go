package core

import "testing"

func TestLedgerCompressionRoundTrip(t *testing.T) {
	l := NewLedger()
	l.Credit("alice", 50)
	b := NewBlock(nil, "")
	b.Hash = "gen" // ensure deterministic
	if err := l.AddBlock(b); err != nil {
		t.Fatalf("add block: %v", err)
	}

	tmp := t.TempDir() + "/snap.gz"
	if err := SaveCompressedSnapshot(l, tmp); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	loaded, err := LoadCompressedSnapshot(tmp)
	if err != nil {
		t.Fatalf("load snapshot: %v", err)
	}
	if h, _ := loaded.Head(); h != 1 {
		t.Fatalf("expected height 1 got %d", h)
	}
	if bal := loaded.GetBalance("alice"); bal != 50 {
		t.Fatalf("unexpected balance %d", bal)
	}
}
