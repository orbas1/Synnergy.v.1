package synnergy

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestIndexingNodeBasic(t *testing.T) {
	ctx := context.Background()
	var events []IndexEvent
	n := NewIndexingNode(WithIndexWatcher(IndexEventHandlerFunc(func(ctx context.Context, event IndexEvent) error {
		events = append(events, event)
		return nil
	})))

	if _, err := n.IndexWithOptions(ctx, "k1", []byte("v1")); err != nil {
		t.Fatalf("index k1: %v", err)
	}
	if _, err := n.IndexWithOptions(ctx, "k2", []byte("v2"), WithEntryTags(map[string]string{"role": "authority"})); err != nil {
		t.Fatalf("index k2: %v", err)
	}
	if n.Count() != 2 {
		t.Fatalf("expected count 2, got %d", n.Count())
	}

	res, ok := n.QueryWithMetadata("k1")
	if !ok || string(res.Value) != "v1" || res.Metadata.Version != 1 {
		t.Fatalf("unexpected result: %+v %v", res, ok)
	}
	res2, ok := n.QueryWithMetadata("k2")
	if !ok || res2.Metadata.Tags["role"] != "authority" {
		t.Fatalf("expected metadata tag, got %+v", res2.Metadata)
	}

	// Ensure snapshot copies data.
	snapshot := n.Snapshot()
	snapshot["k1"] = IndexResult{Value: []byte("tampered")}
	if got, _ := n.Query("k1"); string(got) != "v1" {
		t.Fatalf("snapshot mutation should not affect stored data")
	}

	keys := n.Keys()
	if len(keys) != 2 || keys[0] != "k1" || keys[1] != "k2" {
		t.Fatalf("unexpected keys: %v", keys)
	}

	n.Remove(ctx, "k1")
	if n.Count() != 1 {
		t.Fatalf("expected count 1 after remove, got %d", n.Count())
	}
	if len(events) < 3 { // k1 indexed, k2 indexed, k1 removed at minimum
		t.Fatalf("expected events, got %v", events)
	}
}

func TestIndexingNodeTTLAndEviction(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	var events []IndexEvent
	n := NewIndexingNode(WithIndexClock(clock), WithIndexTTL(2*time.Second), WithIndexMaxEntries(2), WithIndexWatcher(IndexEventHandlerFunc(func(ctx context.Context, event IndexEvent) error {
		events = append(events, event)
		return nil
	})))

	if _, err := n.IndexWithOptions(context.Background(), "k1", []byte("v1")); err != nil {
		t.Fatalf("index k1: %v", err)
	}
	if _, err := n.IndexWithOptions(context.Background(), "k2", []byte("v2")); err != nil {
		t.Fatalf("index k2: %v", err)
	}
	if _, err := n.IndexWithOptions(context.Background(), "k3", []byte("v3")); err != nil {
		t.Fatalf("index k3: %v", err)
	}
	// Since max entries is 2, one key should have been evicted.
	if n.Count() != 2 {
		t.Fatalf("expected eviction to keep 2 entries, got %d", n.Count())
	}

	now = now.Add(3 * time.Second)
	if _, ok := n.Query("k2"); ok {
		t.Fatalf("expected TTL expiration for k2")
	}
	if n.Count() != 1 {
		t.Fatalf("expected one entry after TTL expiration, got %d", n.Count())
	}
	if len(events) < 3 {
		t.Fatalf("expected events emitted for ttl/eviction")
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
			defer wg.Done()
			key := fmt.Sprintf("k%02d", i)
			if _, err := n.IndexWithOptions(context.Background(), key, []byte{byte(i)}); err != nil {
				t.Errorf("index %s: %v", key, err)
			}
		}()
	}
	wg.Wait()
	if n.Count() != ops {
		t.Fatalf("expected %d entries, got %d", ops, n.Count())
	}
	for i := 0; i < ops; i++ {
		key := fmt.Sprintf("k%02d", i)
		res, ok := n.QueryWithMetadata(key)
		if !ok || len(res.Value) != 1 || res.Value[0] != byte(i) {
			t.Fatalf("query mismatch for %s: %v %v", key, res, ok)
		}
	}
}

func TestIndexBatchMergeTags(t *testing.T) {
	n := NewIndexingNode()
	values := map[string][]byte{"a": []byte("1"), "b": []byte("2")}
	if err := n.IndexBatch(context.Background(), values, WithEntryTags(map[string]string{"env": "prod"})); err != nil {
		t.Fatalf("index batch: %v", err)
	}
	if _, err := n.IndexWithOptions(context.Background(), "a", []byte("1b"), WithEntryTags(map[string]string{"zone": "1"}), MergeEntryTags()); err != nil {
		t.Fatalf("merge tags: %v", err)
	}
	res, ok := n.QueryWithMetadata("a")
	if !ok || res.Metadata.Tags["env"] != "prod" || res.Metadata.Tags["zone"] != "1" {
		t.Fatalf("expected merged tags, got %+v", res.Metadata.Tags)
	}
}
