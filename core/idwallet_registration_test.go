package core

import "testing"

func TestIDRegistry(t *testing.T) {
	reg := NewIDRegistry()
	if err := reg.Register("addr1", "metadata"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, ok := reg.Info("addr1"); !ok {
		t.Fatalf("expected wallet registered")
	}
	if err := reg.Register("addr1", "other"); err == nil {
		t.Fatalf("expected error for duplicate registration")
	}
}
