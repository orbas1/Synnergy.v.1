package synnergy

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestSecureStorageStoreRetrieve(t *testing.T) {
	s := NewSecureStorage()
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("rand: %v", err)
	}
	data := []byte("model-bytes")
	if err := s.Store("hash", data, key); err != nil {
		t.Fatalf("store: %v", err)
	}
	out, err := s.Retrieve("hash", key)
	if err != nil {
		t.Fatalf("retrieve: %v", err)
	}
	if !bytes.Equal(out, data) {
		t.Fatalf("got %q want %q", out, data)
	}
}

func TestSecureStorageBadKey(t *testing.T) {
	s := NewSecureStorage()
	if err := s.Store("h", []byte("d"), []byte("short")); err == nil {
		t.Fatalf("expected error for short key")
	}
	if _, err := s.Retrieve("missing", make([]byte, 32)); err == nil {
		t.Fatalf("expected error for missing model")
	}
}
