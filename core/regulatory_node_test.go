package core

import "testing"

func TestRegulatoryNodeSignatureAndLogs(t *testing.T) {
	mgr := NewRegulatoryManager()
	mgr.AddRegulation(Regulation{ID: "r1", MaxAmount: 10})
	node := NewRegulatoryNode("node1", mgr)

	// Unknown wallet should fail
	tx := Transaction{From: "alice", Amount: 5}
	if err := node.ApproveTransaction(tx); err == nil {
		t.Fatalf("expected error for unknown wallet")
	}

	// Register wallet and approve valid transaction
	w, err := NewWallet()
	if err != nil {
		t.Fatalf("wallet: %v", err)
	}
	node.RegisterWallet(w)
	tx.From = w.Address
	if _, err := w.Sign(&tx); err != nil {
		t.Fatalf("sign: %v", err)
	}
	if err := node.ApproveTransaction(tx); err != nil {
		t.Fatalf("expected approval, got %v", err)
	}

	// Amount exceeding regulation should be rejected and logged
	tx.Amount = 20
	if _, err := w.Sign(&tx); err != nil {
		t.Fatalf("sign: %v", err)
	}
	if err := node.ApproveTransaction(tx); err == nil {
		t.Fatalf("expected rejection for excess amount")
	}
	if logs := node.Logs(w.Address); len(logs) != 1 {
		t.Fatalf("expected one log, got %v", logs)
	}

	// Tamper signature
	tx.Amount = 5
	tx.Signature = []byte("bad")
	if err := node.ApproveTransaction(tx); err == nil {
		t.Fatalf("expected invalid signature")
	}
}

func TestFlagRequiresReason(t *testing.T) {
	mgr := NewRegulatoryManager()
	node := NewRegulatoryNode("node", mgr)
	if err := node.FlagEntity("bob", ""); err == nil {
		t.Fatalf("expected error for empty reason")
	}
	if err := node.FlagEntity("bob", "fraud"); err != nil {
		t.Fatalf("flag failed: %v", err)
	}
	if logs := node.Logs("bob"); len(logs) != 1 {
		t.Fatalf("expected one log")
	}
}

func TestApproveTransactionNoManager(t *testing.T) {
	node := NewRegulatoryNode("n", nil)
	tx := Transaction{From: "alice", Amount: 100}
	if err := node.ApproveTransaction(tx); err != nil {
		t.Fatalf("expected approval when manager absent, got %v", err)
	}
}
