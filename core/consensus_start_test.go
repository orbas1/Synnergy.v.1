package core

import (
	"context"
	"testing"
	"time"
)

// TestConsensusServiceStartStop verifies the start-stop lifecycle of the consensus service.
func TestConsensusServiceStartStop(t *testing.T) {
	ledger := NewLedger()
	node := NewNode("n1", "addr", ledger)
	if err := node.SetStake("n1", 1); err != nil {
		t.Fatalf("set stake: %v", err)
	}
	ledger.Mint("alice", 100)
	tx := NewTransaction("alice", "bob", 10, 1, 0)
	if err := node.AddTransaction(tx); err != nil {
		t.Fatalf("add tx: %v", err)
	}
	svc := NewConsensusService(node)
	svc.Start(context.Background(), 10*time.Millisecond)

	// wait up to one second for a block to be mined
	deadline := time.Now().Add(1 * time.Second)
	for time.Now().Before(deadline) {
		if h, _ := svc.Info(); h > 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	svc.Stop()
	height, running := svc.Info()
	if running {
		t.Fatalf("service should be stopped")
	}
	if height == 0 {
		t.Fatalf("expected block to be mined")
	}
}

func TestNextConsensusIntervalAdjustments(t *testing.T) {
	base := 50 * time.Millisecond
	if got := nextConsensusInterval(base, 0, nil, 100); got != base*4 {
		t.Fatalf("expected idle backoff to quadruple base interval, got %v", got)
	}

	busyBlock := &Block{}
	accelerated := nextConsensusInterval(200*time.Millisecond, 500, busyBlock, 100)
	if accelerated != 100*time.Millisecond {
		t.Fatalf("expected busy interval to halve base interval, got %v", accelerated)
	}

	// Ensure min interval is enforced when halving falls below threshold.
	minCheck := nextConsensusInterval(2*time.Millisecond, 10, busyBlock, 1)
	if minCheck < 5*time.Millisecond {
		t.Fatalf("expected minimum interval enforcement, got %v", minCheck)
	}
}
