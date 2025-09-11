package core

import "testing"

func TestPlasmaBridgeOperations(t *testing.T) {
	b := NewPlasmaBridge()
	if err := b.Deposit("alice", "token", 100); err != nil {
		t.Fatalf("deposit: %v", err)
	}
	nonce, err := b.StartExit("alice", "token", 50)
	if err != nil {
		t.Fatalf("start exit: %v", err)
	}
	if err := b.FinalizeExit(nonce); err != nil {
		t.Fatalf("finalize: %v", err)
	}
	ex, err := b.GetExit(nonce)
	if err != nil || !ex.Finalized {
		t.Fatalf("exit not finalized")
	}
	if len(b.ListExits("alice")) != 1 {
		t.Fatalf("list exits failed")
	}

	b.Pause()
	if _, err := b.StartExit("alice", "token", 1); err != ErrBridgePaused {
		t.Fatalf("expected pause error, got %v", err)
	}
}
