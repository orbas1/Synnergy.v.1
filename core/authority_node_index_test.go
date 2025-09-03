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

func TestAuthorityNodeIndexJSON(t *testing.T) {
        idx := NewAuthorityNodeIndex()
        idx.Add(&AuthorityNode{Address: "a1"})
        if _, err := idx.MarshalJSON(); err != nil {
                t.Fatalf("marshal: %v", err)
        }
        snap := idx.Snapshot()
        if len(snap) != 1 {
                t.Fatalf("expected 1 entry, got %d", len(snap))
        }
}
