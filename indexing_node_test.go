package synnergy

import (
	"fmt"
	"sort"
	"sync"
	"testing"
)

// TestIndexingNode_IndexQueryCopies ensures that Index stores a copy of the
// provided value and Query returns a fresh copy on each call.
func TestIndexingNode_IndexQueryCopies(t *testing.T) {
	n := NewIndexingNode()
	key := "alpha"
	val := []byte("data")
	n.Index(key, val)
	// mutate original slice after indexing
	val[0] = 'x'

	got1, ok := n.Query(key)
	if !ok {
		t.Fatalf("expected key %q to exist", key)
	}
	if string(got1) != "data" {
		t.Fatalf("expected stored value to remain 'data', got %q", string(got1))
	}

	// mutate returned slice and ensure underlying value is unaffected
	got1[0] = 'y'
	got2, ok := n.Query(key)
	if !ok {
		t.Fatalf("expected key %q to exist on second query", key)
	}
	if string(got2) != "data" {
		t.Fatalf("value changed after modifying query result, got %q", string(got2))
	}
}

// TestIndexingNode_Remove verifies removing a key deletes its entry.
func TestIndexingNode_Remove(t *testing.T) {
	n := NewIndexingNode()
	n.Index("k", []byte("v"))
	if c := n.Count(); c != 1 {
		t.Fatalf("count = %d, want 1", c)
	}
	n.Remove("k")
	if c := n.Count(); c != 0 {
		t.Fatalf("count after remove = %d, want 0", c)
	}
	if _, ok := n.Query("k"); ok {
		t.Fatalf("expected key to be absent after remove")
	}
}

// TestIndexingNode_KeysAndCount ensures Keys returns all current keys and Count matches.
func TestIndexingNode_KeysAndCount(t *testing.T) {
	n := NewIndexingNode()
	keys := []string{"a", "b", "c"}
	for _, k := range keys {
		n.Index(k, []byte(k))
	}
	got := n.Keys()
	sort.Strings(got)
	if len(got) != len(keys) {
		t.Fatalf("len(Keys) = %d, want %d", len(got), len(keys))
	}
	for i, k := range keys {
		if got[i] != k {
			t.Fatalf("Keys[%d] = %s, want %s", i, got[i], k)
		}
	}
	if c := n.Count(); c != len(keys) {
		t.Fatalf("Count() = %d, want %d", c, len(keys))
	}
}

// TestIndexingNode_ConcurrentIndexing runs many concurrent writes to ensure safety.
func TestIndexingNode_ConcurrentIndexing(t *testing.T) {
	n := NewIndexingNode()
	var wg sync.WaitGroup
	total := 100
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("k%03d", i)
			val := []byte{byte(i)}
			n.Index(key, val)
		}(i)
	}
	wg.Wait()

	if c := n.Count(); c != total {
		t.Fatalf("Count() = %d, want %d", c, total)
	}
}
