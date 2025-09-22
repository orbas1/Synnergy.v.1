package synnergy

import "testing"

func TestRegulatoryNodeDecisions(t *testing.T) {
	mgr := NewRegulatoryManager()
	if err := mgr.AddRegulation(Regulation{ID: "r1", MaxAmount: 10}); err != nil {
		t.Fatalf("add regulation: %v", err)
	}
	node := NewRegulatoryNode("node1", mgr)
	tx := Transaction{From: "alice", Amount: 5}
	if !node.ApproveTransaction(tx) {
		t.Fatalf("expected approval")
	}
	if _, ok := node.LastApproval("alice"); !ok {
		t.Fatalf("expected approval timestamp")
	}

	tx.Amount = 20
	if node.ApproveTransaction(tx) {
		t.Fatalf("expected rejection")
	}
	logs := node.Logs("alice")
	if len(logs) == 0 {
		t.Fatalf("expected flag for alice")
	}
	if logs[0] == "" {
		t.Fatalf("expected descriptive reason")
	}
}

func TestRegulatoryNodeHandlesMissingManager(t *testing.T) {
	node := NewRegulatoryNode("node2", nil)
	if node.ApproveTransaction(Transaction{From: "bob"}) {
		t.Fatalf("approval should fail without manager")
	}
	if len(node.Logs("bob")) == 0 {
		t.Fatalf("expected audit log entry")
	}
}
