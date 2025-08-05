package synnergy

import (
	"bytes"
	"testing"
	"time"
)

func TestContentNetworkNode(t *testing.T) {
	n := NewContentNetworkNode("node1", "addr")
	meta1 := NewContentMeta("cid1", "name1", 4, "hash1")
	meta2 := NewContentMeta("cid2", "name2", 8, "hash2")

	n.Register(meta1)
	n.Register(meta2)

	if got, ok := n.Content(meta1.ID); !ok || got.Name != "name1" {
		t.Fatalf("content1 not registered: %v %v", got, ok)
	}
	if got, ok := n.Content(meta2.ID); !ok || got.Name != "name2" {
		t.Fatalf("content2 not registered: %v %v", got, ok)
	}
	if len(n.List()) != 2 {
		t.Fatalf("expected two items")
	}

	updated := NewContentMeta("cid2", "updated", 9, "hash3")
	n.Register(updated)
	if got, _ := n.Content("cid2"); got.Name != "updated" {
		t.Fatalf("expected updated meta, got %v", got)
	}

	n.Unregister("cid1")
	if _, ok := n.Content("cid1"); ok {
		t.Fatalf("unregister failed")
	}
	if len(n.List()) != 1 {
		t.Fatalf("expected one item after unregister")
	}
}

func TestContentNode(t *testing.T) {
	key := make([]byte, 32)
	node, err := NewContentNode(key)
	if err != nil {
		t.Fatalf("new node: %v", err)
	}
	meta, err := node.StoreContent("greet", []byte("hello"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	if meta.Name != "greet" {
		t.Fatalf("unexpected meta name")
	}
	plain, ok, err := node.RetrieveContent(meta.ID)
	if err != nil || !ok || string(plain) != "hello" {
		t.Fatalf("retrieve mismatch")
	}
	node.DeleteContent(meta.ID)
	if _, ok, _ := node.RetrieveContent(meta.ID); ok {
		t.Fatalf("content not deleted")
	}
}

func TestContentNodeErrorsAndMeta(t *testing.T) {
	if _, err := NewContentNode([]byte("short")); err == nil {
		t.Fatalf("expected error for short key")
	}

	key := make([]byte, 32)
	node, err := NewContentNode(key)
	if err != nil {
		t.Fatalf("new node: %v", err)
	}

	if _, ok := node.Meta("missing"); ok {
		t.Fatalf("expected no meta for missing id")
	}

	meta, err := node.StoreContent("greet", []byte("hello"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	if m, ok := node.Meta(meta.ID); !ok || m.Name != "greet" {
		t.Fatalf("meta mismatch: %v %v", m, ok)
	}
	if _, ok, err := node.RetrieveContent("other"); ok || err != nil {
		t.Fatalf("expected miss without error, got %v %v", ok, err)
	}
}

func TestDataDistribution(t *testing.T) {
	d := NewDataDistribution()
	meta := NewContentMeta("id", "file", 1, "hash")
	d.Offer("n1", meta)
	d.Offer("n1", meta) // duplicate should not add another location
	d.Offer("n2", meta)
	if len(d.Locations("id")) != 2 {
		t.Fatalf("expected two unique locations")
	}
	if m, ok := d.Meta("id"); !ok || m.Name != "file" {
		t.Fatalf("meta lookup failed")
	}
	if _, ok := d.Meta("other"); ok {
		t.Fatalf("expected no meta for unknown id")
	}
	d.Revoke("n1", "id")
	if len(d.Locations("id")) != 1 {
		t.Fatalf("expected one location after revoke")
	}
	d.Revoke("n2", "id")
	if locs := d.Locations("id"); len(locs) != 0 {
		t.Fatalf("expected dataset removed: %v", locs)
	}
}

func TestDataFeed(t *testing.T) {
	f := NewDataFeed("feed")
	f.Update("k1", "v1")
	if v, ok := f.Get("k1"); !ok || v != "v1" {
		t.Fatalf("get returned %v %v", v, ok)
	}
	if len(f.Keys()) != 1 {
		t.Fatalf("expected one key")
	}
	snap := f.Snapshot()
	if snap["k1"] != "v1" {
		t.Fatalf("snapshot mismatch")
	}
	if f.LastUpdated().IsZero() {
		t.Fatalf("last updated not set")
	}
	f.Delete("k1")
	if _, ok := f.Get("k1"); ok {
		t.Fatalf("delete failed")
	}
}

func TestDataFeedLastUpdated(t *testing.T) {
	f := NewDataFeed("feed")
	if !f.LastUpdated().IsZero() {
		t.Fatalf("expected zero timestamp initially")
	}
	f.Update("k", "v")
	first := f.LastUpdated()
	f.Delete("missing")
	if !f.LastUpdated().Equal(first) {
		t.Fatalf("timestamp changed on deleting missing key")
	}
	time.Sleep(time.Millisecond)
	f.Delete("k")
	if !f.LastUpdated().After(first) {
		t.Fatalf("timestamp not updated after delete existing key")
	}
}

func TestDataResourceManager(t *testing.T) {
	m := NewDataResourceManager()
	m.Put("a", []byte{1, 2})
	m.Put("b", []byte{3})
	if m.Usage() != 3 {
		t.Fatalf("unexpected usage: %d", m.Usage())
	}
	if len(m.Keys()) != 2 {
		t.Fatalf("expected two keys")
	}
	if v, ok := m.Get("a"); !ok || !bytes.Equal(v, []byte{1, 2}) {
		t.Fatalf("get mismatch")
	}
	m.Delete("a")
	if _, ok := m.Get("a"); ok {
		t.Fatalf("delete failed")
	}
	if m.Usage() != 1 {
		t.Fatalf("usage not updated")
	}
}

func TestDataResourceManagerIsolation(t *testing.T) {
	m := NewDataResourceManager()
	m.Put("a", []byte{1, 2, 3})
	buf, _ := m.Get("a")
	buf[0] = 9
	if v, _ := m.Get("a"); v[0] == 9 {
		t.Fatalf("internal data mutated by caller")
	}

	m.Put("a", []byte{4})
	if m.Usage() != 1 {
		t.Fatalf("expected usage 1 after overwrite, got %d", m.Usage())
	}
}

func TestIndexingNode(t *testing.T) {
	n := NewIndexingNode()
	n.Index("k", []byte("v"))
	if c := n.Count(); c != 1 {
		t.Fatalf("count %d", c)
	}
	if v, ok := n.Query("k"); !ok || !bytes.Equal(v, []byte("v")) {
		t.Fatalf("query mismatch")
	}
	if len(n.Keys()) != 1 {
		t.Fatalf("keys mismatch")
	}
	n.Remove("k")
	if _, ok := n.Query("k"); ok {
		t.Fatalf("remove failed")
	}
	if n.Count() != 0 {
		t.Fatalf("count mismatch after remove")
	}
}

func TestIndexingNodeEdgeCases(t *testing.T) {
	n := NewIndexingNode()
	if _, ok := n.Query("missing"); ok {
		t.Fatalf("expected missing key to return false")
	}
	n.Remove("missing")
	if n.Count() != 0 {
		t.Fatalf("expected count 0 after removing missing key")
	}
}
