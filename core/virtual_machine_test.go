package core

import (
	"bytes"
	"context"
	"testing"
)

// TestSimpleVM verifies basic start/stop and opcode execution using the
// default light VM profile.
func TestSimpleVM(t *testing.T) {
	vm := NewSimpleVM()
	if vm.Status() {
		t.Fatalf("expected stopped")
	}
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	wasm := []byte{0x00, 0x00, 0x00} // single NOP opcode
	args := []byte{1, 2, 3}
	out, gas, err := vm.Execute(wasm, "", args, 10)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if gas != 1 || !bytes.Equal(out, args) {
		t.Fatalf("unexpected result")
	}
	if err := vm.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
}

// TestVMVariants ensures that the heavy and super light VM profiles operate and
// that the super light profile enforces a strict concurrency limit.
func TestVMVariants(t *testing.T) {
	heavy := NewSimpleVM(VMHeavy)
	super := NewSimpleVM(VMSuperLight)
	_ = heavy.Start()
	_ = super.Start()

	if _, _, err := heavy.Execute([]byte{0, 0, 0}, "", nil, 5); err != nil {
		t.Fatalf("heavy execute: %v", err)
	}

	// occupy the only slot in super light VM
	super.limiter <- struct{}{}
	if _, _, err := super.Execute([]byte{0, 0, 0}, "", nil, 5); err == nil || err.Error() != "vm busy" {
		t.Fatalf("expected busy error from super light VM")
	}
	<-super.limiter
}

func TestVMContextCancel(t *testing.T) {
	vm := NewSimpleVM()
	_ = vm.Start()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := vm.ExecuteContext(ctx, []byte{0, 0, 0}, "", nil, 5); err == nil {
		t.Fatalf("expected cancellation error")
	}
}
