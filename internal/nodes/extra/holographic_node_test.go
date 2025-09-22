package nodes

import (
	"testing"
	"time"

	"synnergy"
)

func TestHolographicNodeStoreRetrieve(t *testing.T) {
	node := NewHolographicNode("h1", WithFrameLimit(2))
	frame := synnergy.HolographicFrame{ID: "f1"}
	node.Store(frame)
	if err := node.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if got, ok := node.Retrieve("f1"); !ok || got.ID != "f1" {
		t.Fatalf("unexpected frame: %#v, ok=%v", got, ok)
	}
	ids := node.Frames()
	if len(ids) != 1 || ids[0] != "f1" {
		t.Fatalf("unexpected ids: %#v", ids)
	}
}

func TestHolographicNodeRetention(t *testing.T) {
	node := NewHolographicNode("h2", WithFrameLimit(1), WithFrameRetention(5*time.Millisecond))
	node.Store(synnergy.HolographicFrame{ID: "a"})
	time.Sleep(10 * time.Millisecond)
	node.Store(synnergy.HolographicFrame{ID: "b"})
	if _, ok := node.Retrieve("a"); ok {
		t.Fatalf("expected frame a to be evicted")
	}
	if _, ok := node.Retrieve("b"); !ok {
		t.Fatalf("expected frame b to exist")
	}
}
