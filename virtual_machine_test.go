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
	opcode := SNVMOpcodes[0]
	code := opcode.Code
	wasm := []byte{byte(code >> 16), byte(code >> 8), byte(code)}
	args := []byte{0xAA, 0xBB}
	out, gasUsed, err := vm.Execute(wasm, "", args, 10)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !bytes.Equal(out, args) {
		t.Fatalf("opcode should echo args")
	}
	expected := GasCost(opcode.Name)
	if expected == 0 {
		expected = DefaultGasCost
	}
	if gasUsed != expected {
		t.Fatalf("unexpected gas used: got %d want %d", gasUsed, expected)
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

func TestSimpleVMCustomGasResolver(t *testing.T) {
	vm := NewSimpleVM()
	vm.SetGasResolver(func(code uint32) (string, uint64) {
		return "custom", 5
	})
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	wasm := []byte{0, 0, 0}
	if _, gasUsed, err := vm.Execute(wasm, "", nil, 10); err != nil {
		t.Fatalf("execute: %v", err)
	} else if gasUsed != 5 {
		t.Fatalf("expected custom gas 5, got %d", gasUsed)
	}
	if _, _, err := vm.Execute(wasm, "", nil, 4); err == nil {
		t.Fatalf("expected gas limit error")
	}
}
