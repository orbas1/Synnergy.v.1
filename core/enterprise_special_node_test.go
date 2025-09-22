package core

import (
	"testing"
	"time"

	"synnergy/internal/nodes"
)

func TestEnterpriseSpecialNodeAttachSnapshot(t *testing.T) {
	agg := NewEnterpriseSpecialNode(nodes.Address("agg-snapshot"))
	if err := agg.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	ledger := NewLedger()
	ledger.Mint("alice", 200)
	node := NewNode("node-1", "addr-1", ledger)

	plugin := EnterpriseNodePluginFromNode("node-1", CombinedRoleConsensus, node, map[string]string{"region": "us-east"})
	if err := agg.AttachPlugin(plugin); err != nil {
		t.Fatalf("attach: %v", err)
	}

	snap := agg.Snapshot()
	if snap.NodeCount != 1 {
		t.Fatalf("expected 1 node, got %d", snap.NodeCount)
	}
	if snap.Roles[CombinedRoleConsensus] != 1 {
		t.Fatalf("expected consensus role count 1, got %d", snap.Roles[CombinedRoleConsensus])
	}
	if len(snap.Plugins) != 1 {
		t.Fatalf("expected 1 plugin summary, got %d", len(snap.Plugins))
	}
	if snap.Plugins[0].Metrics.MempoolSize != 0 {
		t.Fatalf("expected empty mempool, got %d", snap.Plugins[0].Metrics.MempoolSize)
	}
}

func TestEnterpriseSpecialNodeBroadcastAndLedger(t *testing.T) {
	agg := NewEnterpriseSpecialNode(nodes.Address("agg-broadcast"))
	if err := agg.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	ledgerA := NewLedger()
	ledgerB := NewLedger()
	ledgerA.Mint("alice", 150)
	ledgerB.Mint("alice", 90)

	nodeA := NewNode("node-a", "addr-a", ledgerA)
	nodeB := NewNode("node-b", "addr-b", ledgerB)

	if err := agg.AttachPlugin(EnterpriseNodePluginFromNode("node-a", CombinedRoleConsensus, nodeA, nil)); err != nil {
		t.Fatalf("attach node-a: %v", err)
	}
	if err := agg.AttachPlugin(EnterpriseNodePluginFromNode("node-b", CombinedRoleExecution, nodeB, nil)); err != nil {
		t.Fatalf("attach node-b: %v", err)
	}

	tx := NewTransaction("alice", "bob", 10, 1, 0)
	results, err := agg.BroadcastTransaction(tx)
	if err != nil {
		t.Fatalf("broadcast: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 broadcast results, got %d", len(results))
	}
	for id, resErr := range results {
		if resErr != nil {
			t.Fatalf("broadcast error for %s: %v", id, resErr)
		}
	}

	if len(nodeA.Mempool) != 1 {
		t.Fatalf("expected nodeA mempool 1, got %d", len(nodeA.Mempool))
	}
	if len(nodeB.Mempool) != 1 {
		t.Fatalf("expected nodeB mempool 1, got %d", len(nodeB.Mempool))
	}

	if bal := agg.LedgerBalance("alice"); bal != 240 {
		t.Fatalf("expected aggregated balance 240, got %d", bal)
	}
}

func TestEnterpriseSpecialNodeEvents(t *testing.T) {
	agg := NewEnterpriseSpecialNode(nodes.Address("agg-events"))
	if err := agg.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	plugin := EnterpriseNodePlugin{ID: "observer", Role: CombinedRoleAnalytics}
	if err := agg.AttachPlugin(plugin); err != nil {
		t.Fatalf("attach: %v", err)
	}

	watcherID, ch := agg.WatchEvents()
	defer agg.StopWatching(watcherID)

	if err := agg.UpdatePluginLabels("observer", map[string]string{"tier": "gold"}); err != nil {
		t.Fatalf("update labels: %v", err)
	}
	if !agg.DetachPlugin("observer") {
		t.Fatalf("expected observer plugin to be removed")
	}

	deadline := time.After(2 * time.Second)
	received := 0
	for received < 3 {
		select {
		case evt, ok := <-ch:
			if !ok {
				t.Fatalf("watcher channel closed early")
			}
			received++
			if evt.Sequence == 0 {
				t.Fatalf("expected monotonic sequence, got %d", evt.Sequence)
			}
		case <-deadline:
			t.Fatalf("timed out waiting for events; received %d", received)
		}
	}
}
