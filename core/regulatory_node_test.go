package core

import "testing"

func TestRegulatoryNode(t *testing.T) {
	mgr := NewRegulatoryManager()
	mgr.AddRegulation(Regulation{ID: "r1", MaxAmount: 10})
	node := NewRegulatoryNode("node1", mgr)
	tx := Transaction{From: "alice", Amount: 5}
	if err := node.ApproveTransaction(tx); err != nil {
		t.Fatalf("expected approval: %v", err)
	}
	tx.Amount = 20
	if err := node.ApproveTransaction(tx); err == nil {
		t.Fatalf("expected rejection")
	}
	logs := node.Logs("alice")
	if len(logs) != 1 {
		t.Fatalf("expected flag for alice")
	}
	node.ClearLogs("alice")
	if len(node.Logs("alice")) != 0 {
		t.Fatalf("expected logs cleared")
	}
}
