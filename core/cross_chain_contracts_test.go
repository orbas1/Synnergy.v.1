package core

import "testing"

func TestCrossChainRegistry(t *testing.T) {
	reg := NewCrossChainRegistry()
	reg.RegisterMapping("local1", "chainB", "remote1")
	if _, ok := reg.GetMapping("local1"); !ok {
		t.Fatalf("mapping not found")
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
