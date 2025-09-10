package core

import "testing"

func TestCrossChainRegistry(t *testing.T) {
	reg := NewCrossChainRegistry()

	if err := reg.RegisterMapping("relayer1", "local1", "chainB", "remote1"); err == nil {
		t.Fatalf("expected unauthorized relayer error")
	}

	reg.AuthorizeRelayer("relayer1")
	if err := reg.RegisterMapping("relayer1", "local1", "chainB", "remote1"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if _, ok := reg.GetMapping("local1"); !ok {
		t.Fatalf("mapping not found")
	}
	if len(reg.ListMappings()) != 1 {
		t.Fatalf("expected one mapping")
	}
	if err := reg.RemoveMapping("bad", "local1"); err == nil {
		t.Fatalf("expected unauthorized removal error")
	}
	if err := reg.RemoveMapping("relayer1", "local1"); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if _, ok := reg.GetMapping("local1"); ok {
		t.Fatalf("expected mapping to be removed")
	}
}
