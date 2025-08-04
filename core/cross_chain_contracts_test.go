package core

import "testing"

func TestCrossChainContractRegistry(t *testing.T) {
        reg := NewCrossChainContractRegistry()
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
