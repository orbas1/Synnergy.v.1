package security

import "testing"

func TestKeyManager(t *testing.T) {
	km := NewKeyManager(1)
	if km.Key() != 1 {
		t.Fatalf("expected key 1, got %d", km.Key())
	}
	km.Rotate(2)
	if km.Key() != 2 {
		t.Fatalf("expected key 2, got %d", km.Key())
	}
}
