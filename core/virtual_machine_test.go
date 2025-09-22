package core

import (
	"bytes"
	"context"
	"testing"
	"time"

	synn "synnergy"
)

func TestSimpleVM(t *testing.T) {
	vm := NewSimpleVM()
	defer vm.Close()
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
	if gas == 0 || !bytes.Equal(out, args) {
		t.Fatalf("unexpected result")
	}
	if err := vm.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
}

func TestVMVariants(t *testing.T) {
	heavy := NewSimpleVM(VMHeavy)
	defer heavy.Close()
	super := NewSimpleVM(VMSuperLight)
	defer super.Close()
	_ = heavy.Start()
	_ = super.Start()

	if _, _, err := heavy.Execute([]byte{0, 0, 0}, "", nil, 5); err != nil {
		t.Fatalf("heavy execute: %v", err)
	}

	super.limiter <- struct{}{}
	if _, _, err := super.Execute([]byte{0, 0, 0}, "", nil, 5); err == nil || err.Error() != "vm busy" {
		t.Fatalf("expected busy error from super light VM")
	}
	<-super.limiter
}

func TestVMContextCancel(t *testing.T) {
	vm := NewSimpleVM()
	defer vm.Close()
	_ = vm.Start()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := vm.ExecuteContext(ctx, []byte{0, 0, 0}, "", nil, 5); err == nil {
		t.Fatalf("expected cancellation error")
	}
}

func TestVMCustomHandler(t *testing.T) {
	vm := NewSimpleVM()
	defer vm.Close()
	_ = vm.Start()
	vm.RegisterHandler(0xFFFFFF, func(in []byte) ([]byte, error) { return append(in, 0xAA), nil })
	out, _, err := vm.Execute([]byte{0xFF, 0xFF, 0xFF}, "", []byte{0x01}, 10)
	if err != nil {
		t.Fatalf("exec: %v", err)
	}
	if len(out) != 2 || out[1] != 0xAA {
		t.Fatalf("handler not invoked")
	}
}

func TestVMRespectsGasTableUpdates(t *testing.T) {
	synn.ResetGasTable()
	if err := synn.RegisterGasCost("core_vm_test", 7); err != nil {
		t.Fatalf("register gas: %v", err)
	}
	vm := NewSimpleVM()
	defer vm.Close()
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	vm.RegisterHandlerNamed(0xAA5501, "core_vm_test", func(in []byte) ([]byte, error) { return in, nil })
	if _, _, err := vm.Execute([]byte{0xAA, 0x55, 0x01}, "", nil, 6); err == nil {
		t.Fatalf("expected gas limit error")
	}
	if _, gas, err := vm.Execute([]byte{0xAA, 0x55, 0x01}, "", nil, 7); err != nil {
		t.Fatalf("execute: %v", err)
	} else if gas != 7 {
		t.Fatalf("expected gas usage 7, got %d", gas)
	}
}

func TestVMHotReloadGas(t *testing.T) {
	synn.ResetGasTable()
	if err := synn.RegisterGasCost("core_vm_reload", 2); err != nil {
		t.Fatalf("register: %v", err)
	}
	vm := NewSimpleVM()
	defer vm.Close()
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	vm.RegisterHandlerNamed(0xDDCCEE, "core_vm_reload", func(in []byte) ([]byte, error) { return in, nil })
	if _, _, err := vm.Execute([]byte{0xDD, 0xCC, 0xEE}, "", nil, 2); err != nil {
		t.Fatalf("baseline execute: %v", err)
	}
	if err := synn.RegisterGasCost("core_vm_reload", 11); err != nil {
		t.Fatalf("update gas: %v", err)
	}
	waitForCoreGasCost(t, vm, "core_vm_reload", 11)
	if _, _, err := vm.Execute([]byte{0xDD, 0xCC, 0xEE}, "", nil, 5); err == nil {
		t.Fatalf("expected gas limit exceeded after reload")
	}
}

func waitForCoreGasCost(t *testing.T, vm *SimpleVM, opcode string, expected uint64) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		snap := vm.GasSnapshot()
		if snap.Table[opcode] == expected {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("gas update for %s not observed", opcode)
}
