package core

import "testing"

func TestShardManager(t *testing.T) {
	m := NewShardManager(2) // 4 shards
	m.SetLeader(0, "leader0")
	if l, ok := m.GetLeader(0); !ok || l != "leader0" {
		t.Fatalf("get leader")
	}
	m.SubmitCrossShardTx(0, 1, "tx1")
	m.SubmitCrossShardTx(0, 1, "tx2")
	receipts := m.PullReceipts(1)
	if len(receipts) != 2 {
		t.Fatalf("pull receipts")
	}
	m.Reshard(3)
	if m.ShardCount() != 8 {
		t.Fatalf("shard count")
	}
	m.SubmitCrossShardTx(0, 2, "tx3")
	heavy := m.Rebalance(0)
	if len(heavy) == 0 {
		t.Fatalf("expected heavy shard")
	}

	if load := m.ShardLoad(2); load != 1 {
		t.Fatalf("unexpected shard load: %d", load)
	}
	loads := m.LoadMap()
	if loads[2] != 1 {
		t.Fatalf("load map mismatch")
	}
}
