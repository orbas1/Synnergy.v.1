package core

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestMiningNode(t *testing.T) {
	mn := NewMiningNode(100)
	if mn.IsMining() {
		t.Fatalf("expected node to be stopped initially")
	}
	mn.Start()
	if !mn.IsMining() {
		t.Fatalf("node should be mining after Start")
	}
	hash, err := mn.Mine([]byte("data"))
	if err != nil || hash == "" {
		t.Fatalf("mine returned invalid hash: %v %s", err, hash)
	}
	mn.Stop()
	if mn.IsMining() {
		t.Fatalf("node should be stopped after Stop")
	}
	if _, err := mn.Mine([]byte("data")); err == nil {
		t.Fatalf("expected error when mining is inactive")
	}
}

func TestMiningNodeMineUntil(t *testing.T) {
	mn := NewMiningNode(10)
	mn.Start()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	hash, _, err := mn.MineUntil(ctx, []byte("data"), "0")
	if err != nil {
		t.Fatalf("MineUntil failed: %v", err)
	}
	if !strings.HasPrefix(hash, "0") {
		t.Fatalf("hash %s does not have expected prefix", hash)
	}
	mn.Stop()
}
