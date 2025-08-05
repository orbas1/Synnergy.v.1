package synnergy

import (
	"fmt"
	"sync"
	"testing"
)

func TestIndexingNodeBasic(t *testing.T) {
	n := NewIndexingNode()
	if n.Count() != 0 {
		t.Fatalf("expected empty index, got %d", n.Count())
	}
	n.Index("k1", []byte("v1"))
	n.Index("k2", []byte("v2"))
	if n.Count() != 2 {
		t.Fatalf("expected count 2, got %d", n.Count())
	}
	if v, ok := n.Query("k1"); !ok || string(v) != "v1" {
		t.Fatalf("query returned %q %v", v, ok)
	}
	n.Index("k1", []byte("v1-updated"))
	if v, ok := n.Query("k1"); !ok || string(v) != "v1-updated" {
		t.Fatalf("update failed: %q %v", v, ok)
	}
	keys := n.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	seen := make(map[string]bool)
	for _, k := range keys {
		seen[k] = true
	}
	if !seen["k1"] || !seen["k2"] {
		t.Fatalf("unexpected keys: %v", keys)
	}
	n.Remove("k1")
	if _, ok := n.Query("k1"); ok {
		t.Fatalf("remove failed")
	}
	if n.Count() != 1 {
		t.Fatalf("expected count 1, got %d", n.Count())
	}
}

func TestIndexingNodeQueryReturnsCopy(t *testing.T) {
	n := NewIndexingNode()
	orig := []byte("data")
	n.Index("k", orig)
	orig[0] = 'x'
	v, ok := n.Query("k")
	if !ok || string(v) != "data" {
		t.Fatalf("expected stored copy unaffected: %q %v", v, ok)
	}
	v[0] = 'z'
	v2, _ := n.Query("k")
	if string(v2) != "data" {
		t.Fatalf("query should return copy, got %q", v2)
	}
}

func TestIndexingNodeConcurrentAccess(t *testing.T) {
	n := NewIndexingNode()
	const ops = 100
	var wg sync.WaitGroup
	for i := 0; i < ops; i++ {
		i := i
		wg.Add(1)
		go func() {
			n.Index(fmt.Sprintf("k%02d", i), []byte{byte(i)})
			wg.Done()
		}()
	}
	wg.Wait()
	if n.Count() != ops {
		t.Fatalf("expected %d entries, got %d", ops, n.Count())
	}
	for i := 0; i < ops; i++ {
		key := fmt.Sprintf("k%02d", i)
		if v, ok := n.Query(key); !ok || len(v) != 1 || v[0] != byte(i) {
			t.Fatalf("query mismatch for %s: %v %v", key, v, ok)
		}
	}
}
