package core

import "testing"

func TestContractRegistry(t *testing.T) {
	vm := NewSimpleVM()
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	reg := NewContractRegistry(vm)
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
}
