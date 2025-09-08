package security

import "testing"

func TestSecretsManager(t *testing.T) {
	sm := NewSecretsManager()
	if err := sm.Store("k", "v"); err != nil {
		t.Fatalf("store: %v", err)
	}
	if err := sm.Store("", ""); err == nil {
		t.Fatalf("expected error for empty key/value")
	}
	val, err := sm.Retrieve("k")
	if err != nil || val != "v" {
		t.Fatalf("expected v, got %s (%v)", val, err)
	}
	if _, err := sm.Retrieve("missing"); err == nil {
		t.Fatalf("expected error for missing key")
	}
}
