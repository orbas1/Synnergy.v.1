package core

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// TestSimpleVM verifies basic start/stop and opcode execution using the
// default light VM profile.
func TestSimpleVM(t *testing.T) {
	vm := NewSimpleVM()
	vm.ResetMetrics()
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

	metrics := vm.Metrics()
	if metrics.Executions != 1 {
		t.Fatalf("expected 1 execution, got %d", metrics.Executions)
	}
	if metrics.GasConsumed != 1 {
		t.Fatalf("expected gas consumed 1, got %d", metrics.GasConsumed)
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

// TestVMCustomHandler ensures dynamically registered opcode handlers execute
// correctly and override existing entries.
func TestVMCustomHandler(t *testing.T) {
	vm := NewSimpleVM()
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

func TestVMGasLimitEnforced(t *testing.T) {
	vm := NewSimpleVM()
	_ = vm.Start()
	_, _, err := vm.Execute([]byte{0, 0, 0, 0, 0, 0}, "", nil, 1)
	if err == nil {
		t.Fatalf("expected gas limit error")
	}
	if !errors.Is(err, ErrGasLimit) {
		t.Fatalf("expected ErrGasLimit, got %v", err)
	}
}

func TestVMHooksCaptureTraces(t *testing.T) {
	vm := NewSimpleVM()
	_ = vm.Start()
	defer vm.ResetHooks()

	var (
		mu     sync.Mutex
		traces []ExecutionTrace
	)
	vm.RegisterHook(func(trace ExecutionTrace) {
		mu.Lock()
		defer mu.Unlock()
		traces = append(traces, trace)
	})

	wasm := []byte{0, 0, 0, 0, 0, 0} // two opcodes
	args := []byte{0x01}
	out, gas, err := vm.Execute(wasm, "", args, 10)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !bytes.Equal(out, args) {
		t.Fatalf("expected echo output")
	}
	if gas != 2 {
		t.Fatalf("expected gas 2, got %d", gas)
	}

	time.Sleep(5 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(traces) != 2 {
		t.Fatalf("expected 2 traces, got %d", len(traces))
	}
	for _, trace := range traces {
		if trace.Name == "" {
			t.Fatalf("trace missing opcode name")
		}
		if trace.GasCost == 0 {
			t.Fatalf("trace missing gas cost")
		}
	}
}

func TestSimpleVMCallMeterUnlimited(t *testing.T) {
	vm := NewSimpleVM()
	for i := 0; i < 128; i++ {
		if err := vm.Gas(1); err != nil {
			t.Fatalf("unexpected gas error: %v", err)
		}
	}
	if vm.CallGasLimit() != 0 {
		t.Fatalf("expected unlimited gas meter")
	}
	if vm.CallGasUsed() != 128 {
		t.Fatalf("expected used gas to reflect consumed charges")
	}
}

func TestSimpleVMCallMeterLimitAndRefill(t *testing.T) {
	vm := NewSimpleVM()
	vm.ConfigureCallMeter(5, time.Millisecond)

	for i := 0; i < 5; i++ {
		if err := vm.Gas(1); err != nil {
			t.Fatalf("unexpected gas error: %v", err)
		}
	}
	if err := vm.Gas(1); err == nil {
		t.Fatalf("expected gas limit error once meter is exhausted")
	} else if !errors.Is(err, ErrGasLimit) {
		t.Fatalf("expected ErrGasLimit, got %v", err)
	}

	time.Sleep(2 * time.Millisecond)
	if remaining := vm.CallGasRemaining(); remaining != 5 {
		t.Fatalf("expected meter to refill, remaining %d", remaining)
	}
	if err := vm.Gas(3); err != nil {
		t.Fatalf("unexpected error after refill: %v", err)
	}
	vm.DisableCallMeter()
	if vm.CallGasLimit() != 0 || vm.CallGasRemaining() != 0 {
		t.Fatalf("expected disabled meter to report unlimited usage")
	}
}
