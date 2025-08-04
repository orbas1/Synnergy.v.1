package core

import "testing"

func TestContractRegistry(t *testing.T) {
	r := NewContractRegistry()
	r.RegisterMapping("loc1", "chainB", "rem1")
	m, err := r.GetMapping("loc1")
	if err != nil || m.RemoteAddress != "rem1" {
		t.Fatalf("unexpected mapping: %#v err=%v", m, err)
	}
	if len(r.ListMappings()) != 1 {
		t.Fatalf("expected one mapping")
	}
	if err := r.RemoveMapping("loc1"); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if _, err := r.GetMapping("loc1"); err == nil {
		t.Fatalf("expected error after removal")
	}
}
