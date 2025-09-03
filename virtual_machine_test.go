package synnergy

import (
	"bytes"
	"context"
	"testing"
)

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

func TestSNVMOpcodeExecution(t *testing.T) {
	if len(SNVMOpcodes) == 0 {
		t.Fatalf("no opcodes defined")
	}
	vm := NewSimpleVM()
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	code := SNVMOpcodes[0].Code
	wasm := []byte{byte(code >> 16), byte(code >> 8), byte(code)}
	args := []byte{0xAA, 0xBB}
	out, _, err := vm.Execute(wasm, "", args, 10)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !bytes.Equal(out, args) {
		t.Fatalf("opcode should echo args")
	}
}

func TestVMContextCancel(t *testing.T) {
	vm := NewSimpleVM()
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, err := vm.ExecuteContext(ctx, []byte{0, 0, 0}, "", nil, 10)
	if err == nil {
		t.Fatalf("expected context cancellation error")
	}
}
