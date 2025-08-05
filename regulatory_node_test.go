package synnergy

import "testing"

func TestRegulatoryNode(t *testing.T) {
	mgr := NewRegulatoryManager()
	mgr.AddRegulation(Regulation{ID: "r1", MaxAmount: 10})
	node := NewRegulatoryNode("node1", mgr)
	tx := Transaction{Amount: 5}
	if !node.ApproveTransaction(tx) {
		t.Fatalf("expected approval")
	}
	tx.Amount = 20
	if node.ApproveTransaction(tx) {
		t.Fatalf("expected rejection")
	}
	node.FlagEntity("addr", "suspicious")
	logs := node.Logs("addr")
	if len(logs) != 1 || logs[0] != "suspicious" {
		t.Fatalf("log not recorded")
	}
}
