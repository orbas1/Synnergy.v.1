package synnergy

import "testing"

func TestIDRegistry(t *testing.T) {
	reg := NewIDRegistry()
	if err := reg.Register("addr1", "info"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := reg.Register("addr1", "info"); err == nil {
		t.Fatalf("expected duplicate registration error")
	}
	if !reg.IsRegistered("addr1") {
		t.Fatalf("address should be registered")
	}
	info, ok := reg.Info("addr1")
	if !ok || info != "info" {
		t.Fatalf("unexpected info: %v %v", info, ok)
	}
}
