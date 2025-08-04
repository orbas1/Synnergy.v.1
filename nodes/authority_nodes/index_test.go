package authority_nodes

import "testing"

func TestIndex(t *testing.T) {
	idx := NewIndex()
	n := &AuthorityNode{Address: "a", Role: "r", Votes: make(map[string]bool)}
	idx.Add(n)
	if got, ok := idx.Get("a"); !ok || got.Address != "a" {
		t.Fatalf("get failed")
	}
	if len(idx.List()) != 1 {
		t.Fatalf("list size incorrect")
	}
	idx.Remove("a")
	if _, ok := idx.Get("a"); ok {
		t.Fatalf("remove failed")
	}
}
