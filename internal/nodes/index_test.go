package nodes

import "testing"

func TestBasicNodeLifecycle(t *testing.T) {
	n := NewBasicNode("n1")
	if n.IsRunning() {
		t.Fatal("expected not running")
	}
	if err := n.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if !n.IsRunning() {
		t.Fatal("expected running")
	}
	if err := n.DialSeed("peer1"); err != nil {
		t.Fatalf("dial: %v", err)
	}
	peers := n.Peers()
	if len(peers) != 1 || peers[0] != "peer1" {
		t.Fatalf("unexpected peers: %v", peers)
	}
	if err := n.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if n.IsRunning() {
		t.Fatal("expected stopped")
	}
}
