package security

import "testing"

func TestSecretsManager(t *testing.T) {
	sm := NewSecretsManager()
	sm.Store("k", "v")
	if val, ok := sm.Retrieve("k"); !ok || val != "v" {
		t.Fatalf("expected v, got %s", val)
	}
}
