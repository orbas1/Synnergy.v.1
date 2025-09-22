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
	if err := heavy.Start(); err != nil {
		t.Fatalf("heavy start: %v", err)
	}
	defer heavy.Stop()
	if err := super.Start(); err != nil {
		t.Fatalf("super start: %v", err)
	}
	defer super.Stop()

	if _, _, err := heavy.Execute([]byte{0, 0, 0}, "", nil, 5); err != nil {
		t.Fatalf("heavy execute: %v", err)
	}

	slowOpcode := uint32(0x010203)
	started := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	super.RegisterHandler(slowOpcode, func(in []byte) ([]byte, error) {
		once.Do(func() { close(started) })
		<-release
		return in, nil
	})

	bytecode := []byte{0x01, 0x02, 0x03}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, _, err := super.Execute(bytecode, "", nil, 5); err != nil {
			t.Errorf("first execution failed: %v", err)
		}
	}()

	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, _, err := super.ExecuteContext(ctx, bytecode, "", nil, 5)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}

	close(release)
	wg.Wait()
}

func TestVMStopCancelsQueuedExecutions(t *testing.T) {
	vm := NewSimpleVM(VMSuperLight)
	slowOpcode := uint32(0x020304)
	started := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	vm.RegisterHandler(slowOpcode, func(in []byte) ([]byte, error) {
		once.Do(func() { close(started) })
		<-release
		return in, nil
	})

	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	bytecode := []byte{0x02, 0x03, 0x04}
	firstDone := make(chan error, 1)
	go func() {
		_, _, err := vm.Execute(bytecode, "", nil, 5)
		firstDone <- err
	}()

	<-started

	secondResult := make(chan error, 1)
	go func() {
		_, _, err := vm.Execute(bytecode, "", nil, 5)
		secondResult <- err
	}()

	time.Sleep(5 * time.Millisecond)

	stopErr := make(chan error, 1)
	go func() { stopErr <- vm.Stop() }()

	select {
	case err := <-secondResult:
		if !errors.Is(err, errVMNotRunning) {
			t.Fatalf("expected errVMNotRunning, got %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timed out waiting for queued execution to abort")
	}

	close(release)
	if err := <-firstDone; err != nil {
		t.Fatalf("first execution failed: %v", err)
	}
	if err := <-stopErr; err != nil {
		t.Fatalf("stop failed: %v", err)
	}
}

func TestVMContextCancel(t *testing.T) {
	vm := NewSimpleVM()
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer vm.Stop()
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
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer vm.Stop()
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
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer vm.Stop()
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
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer vm.Stop()
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
