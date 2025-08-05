package synnergy

import (
	"strings"
	"testing"
	"time"
)

func TestMiningNodeLifecycle(t *testing.T) {
	m := NewMiningNode("miner1", 1)
	if m.ID() != "miner1" {
		t.Fatalf("ID mismatch: %s", m.ID())
	}
	if m.HashRate() != 1 {
		t.Fatalf("hash rate mismatch")
	}
	m.Start()
	if !m.IsRunning() {
		t.Fatalf("expected running after Start")
	}
	time.Sleep(1100 * time.Millisecond)
	if lb := m.LastBlock(); lb == "" {
		t.Fatalf("expected mined block hash")
	}
	m.Stop()
	if m.IsRunning() {
		t.Fatalf("expected not running after Stop")
	}
}

func TestMiningNodeMineBlock(t *testing.T) {
	m := NewMiningNode("miner", 0)
	hash, err := m.MineBlock(4)
	if err != nil {
		t.Fatalf("mine: %v", err)
	}
	if !strings.HasPrefix(hash, "0") {
		t.Fatalf("hash does not meet difficulty: %s", hash)
	}
	if _, err := m.MineBlock(-1); err == nil {
		t.Fatalf("expected error for negative difficulty")
	}
}

func TestMiningNodeSubmitBlock(t *testing.T) {
	m := NewMiningNode("miner", 0)
	m.SubmitBlock("abc")
	if lb := m.LastBlock(); lb != "abc" {
		t.Fatalf("last block mismatch: %s", lb)
	}
}
