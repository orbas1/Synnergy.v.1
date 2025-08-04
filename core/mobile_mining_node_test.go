package core

import "testing"

func TestMobileMiningNode(t *testing.T) {
	mm := NewMobileMiningNode(50, 1)
	mm.Start()
	if !mm.IsMining() {
		t.Fatalf("expected mining to be active")
	}
	if _, err := mm.Mine([]byte("block")); err != nil {
		t.Fatalf("mine failed: %v", err)
	}
	mm.SetPowerLimit(0)
	if hash, _ := mm.Mine([]byte("block")); hash != "" {
		t.Fatalf("expected empty hash when power limit is zero")
	}
	mm.Stop()
	if mm.IsMining() {
		t.Fatalf("expected mining to stop")
	}
}
