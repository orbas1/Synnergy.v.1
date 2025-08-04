package synnergy

import (
	"bytes"
	"testing"
)

func TestContentNetworkNode(t *testing.T) {
	n := NewContentNetworkNode("node1", "addr")
	meta := NewContentMeta("cid", "name", 4, "hash")
	n.Register(meta)
	if got, ok := n.Content("cid"); !ok || got.Name != "name" {
		t.Fatalf("content not registered: %v %v", got, ok)
	}
	if len(n.List()) != 1 {
		t.Fatalf("expected one item")
	}
	n.Unregister("cid")
	if _, ok := n.Content("cid"); ok {
		t.Fatalf("unregister failed")
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

func TestDataDistribution(t *testing.T) {
	d := NewDataDistribution()
	meta := NewContentMeta("id", "file", 1, "hash")
	d.Offer("n1", meta)
	d.Offer("n2", meta)
	if len(d.Locations("id")) != 2 {
		t.Fatalf("expected two locations")
	}
	if m, ok := d.Meta("id"); !ok || m.Name != "file" {
		t.Fatalf("meta lookup failed")
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
