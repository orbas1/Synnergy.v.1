package synnergy

import "testing"

func TestSimpleVMLifecycle(t *testing.T) {
	vm := NewSimpleVM()
	if vm.Status() {
		t.Fatalf("expected stopped VM")
	}
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if !vm.Status() {
		t.Fatalf("expected running VM")
	}
	if err := vm.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if vm.Status() {
		t.Fatalf("expected stopped VM")
	}
}
