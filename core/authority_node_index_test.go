package core

import "testing"

func TestAuthorityNodeIndex(t *testing.T) {
	idx := NewAuthorityNodeIndex()
	node := &AuthorityNode{Address: "a1"}
	idx.Add(node)
	if _, ok := idx.Get("a1"); !ok {
		t.Fatalf("node not found")
	}
	idx.Remove("a1")
	if _, ok := idx.Get("a1"); ok {
		t.Fatalf("node should be removed")
	}
}
