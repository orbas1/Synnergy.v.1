package core

import "testing"

// TestSwarmBroadcast ensures transactions are broadcast to all swarm members.
func TestSwarmBroadcast(t *testing.T) {
	ledger := NewLedger()
	ledger.Credit("alice", 100)

	n1 := NewNode("n1", "addr1", ledger)
	n2 := NewNode("n2", "addr2", ledger)

	s := NewSwarm()
	s.Join(n1)
	s.Join(n2)

	tx := NewTransaction("alice", "bob", 1, 1, 1)
	s.Broadcast(tx)

	if len(n1.Mempool) != 1 || len(n2.Mempool) != 1 {
		t.Fatalf("transaction was not broadcast to all nodes")
	}
}
