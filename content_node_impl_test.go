package synnergy

import "testing"

func TestContentNode_StoreRetrieve(t *testing.T) {
	key := make([]byte, 32)
	n, err := NewContentNode(key)
	if err != nil {
		t.Fatalf("new node: %v", err)
	}
	meta, err := n.StoreContent("file", []byte("data"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	got, ok, err := n.RetrieveContent(meta.ID)
	if err != nil || !ok || string(got) != "data" {
		t.Fatalf("retrieve mismatch: %v %v %v", string(got), ok, err)
	}
	n.DeleteContent(meta.ID)
	if _, ok, _ := n.RetrieveContent(meta.ID); ok {
		t.Fatalf("content should be removed")
	}
}

func TestContentNode_KeyLength(t *testing.T) {
	if _, err := NewContentNode([]byte("short")); err == nil {
		t.Fatalf("expected error for short key")
	}
}
