package core

import "testing"

func TestSimpleVM(t *testing.T) {
	vm := NewSimpleVM()
	if vm.Status() {
		t.Fatalf("expected stopped")
	}
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	out, gas, err := vm.Execute(nil, "", []byte{1, 2, 3}, 10)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if gas != 3 || len(out) != 3 {
		t.Fatalf("unexpected result")
	}
	if err := vm.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
