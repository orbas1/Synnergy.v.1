package core

import "testing"

func TestSidechainRegistry(t *testing.T) {
	reg := NewSidechainRegistry()
	if _, err := reg.Register("chain1", "meta", []string{"val1"}); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := reg.SubmitHeader("chain1", "h1"); err != nil {
		t.Fatalf("submit header: %v", err)
	}
	if h, ok := reg.GetHeader("chain1"); !ok || h != "h1" {
		t.Fatalf("get header failed")
	}
	if err := reg.Pause("chain1"); err != nil {
		t.Fatalf("pause: %v", err)
	}
	if err := reg.Resume("chain1"); err != nil {
		t.Fatalf("resume: %v", err)
	}
	if err := reg.UpdateValidators("chain1", []string{"val2"}); err != nil {
		t.Fatalf("update validators: %v", err)
	}
	if reg.chains["chain1"].Validators[0] != "val2" {
		t.Fatalf("validators not updated")
	}
	if len(reg.List()) != 1 {
		t.Fatalf("list length")
	}
	if err := reg.Remove("chain1"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if len(reg.List()) != 0 {
		t.Fatalf("remove failed")
	}
}
