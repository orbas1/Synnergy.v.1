package core

import "testing"

func TestContractRegistry(t *testing.T) {
	reg := NewContractRegistry()
	reg.RegisterMapping("local1", "chainB", "remote1")
	if _, ok := reg.GetMapping("local1"); !ok {
		t.Fatalf("mapping not found")
	}
	if len(reg.ListMappings()) != 1 {
		t.Fatalf("list: expected 1 mapping")
	}
	reg.RemoveMapping("local1")
	if _, ok := reg.GetMapping("local1"); ok {
		t.Fatalf("mapping should be removed")
	}
}
