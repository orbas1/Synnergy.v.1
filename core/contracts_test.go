package core

import "testing"

func TestContractRegistry(t *testing.T) {
	vm := NewSimpleVM()
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	ledger := NewLedger()
	ledger.Credit("owner", 1_000)
	reg := NewContractRegistry(vm, ledger)
	wasm := []byte{0x00, 0x61}
	addr, err := reg.Deploy(wasm, "", 10, "owner")
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	if len(reg.List()) != 1 {
		t.Fatalf("expected 1 contract")
	}
	out, _, err := reg.Invoke(addr, "echo", []byte("hi"), 10)
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	if string(out) != "hi" {
		t.Fatalf("unexpected output: %q", out)
	}
	if _, err := reg.Deploy(wasm, "", 10, "owner"); err == nil {
		t.Fatalf("expected duplicate deploy error")
	}
	if _, ok := reg.Get("missing"); ok {
		t.Fatalf("expected missing contract")
	}
}
