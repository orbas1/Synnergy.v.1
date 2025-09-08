package core

import "testing"

func TestContentNetworkNode_RegisterUnregister(t *testing.T) {
	n := NewContentNetworkNode("node1", "addr")
	meta := NewContentMeta("cid", "name", 4, "hash")
	if err := n.Register(meta); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := n.Register(meta); err == nil {
		t.Fatalf("expected duplicate register error")
	}
	if err := n.Unregister(meta.ID); err != nil {
		t.Fatalf("unregister: %v", err)
	}
	if err := n.Unregister(meta.ID); err == nil {
		t.Fatalf("expected missing unregister error")
	}
}

func TestContentNetworkNode_ContentAndList(t *testing.T) {
	n := NewContentNetworkNode("node1", "addr")
	meta := NewContentMeta("cid", "name", 4, "hash")
	if err := n.Register(meta); err != nil {
		t.Fatalf("register: %v", err)
	}
	got, ok := n.Content(meta.ID)
	if !ok || got.Name != "name" {
		t.Fatalf("content lookup failed: %v %v", got, ok)
	}
	list := n.List()
	if len(list) != 1 || list[0].ID != meta.ID {
		t.Fatalf("list mismatch: %v", list)
	}
}
