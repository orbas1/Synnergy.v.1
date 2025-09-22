package synnergy

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"testing"
	"time"
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
	defer vm.Stop()
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
	defer vm.Stop()
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
	defer vm.Stop()
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

func TestSimpleVMAllowsPartialOpcodePadding(t *testing.T) {
	vm := NewSimpleVM()
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer vm.Stop()

	observed := make(chan ExecutionStats, 1)
	vm.SetObserver(func(stats ExecutionStats) {
		observed <- stats
	})

	if out, gasUsed, err := vm.Execute([]byte{0xAA}, "", []byte{0x01}, 10); err != nil {
		t.Fatalf("execute: %v", err)
	} else if len(out) != 1 || out[0] != 0x01 {
		t.Fatalf("unexpected output: %v", out)
	} else if gasUsed == 0 {
		t.Fatalf("expected gas usage")
	}

	select {
	case stats := <-observed:
		if stats.PaddedOpcodes == 0 {
			t.Fatalf("expected padding to be recorded")
		}
	case <-time.After(time.Second):
		t.Fatalf("observer not invoked")
	}
}

func TestSimpleVMBusyWaitsForSlot(t *testing.T) {
	vm := NewSimpleVM(VMSuperLight)
	slowOpcode := uint32(0x010203)
	started := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	vm.RegisterOpcode(slowOpcode, func(in []byte) ([]byte, error) {
		once.Do(func() {
			close(started)
			<-release
		})
		return in, nil
	})
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer vm.Stop()

	bytecode := []byte{0x01, 0x02, 0x03}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, _, err := vm.Execute(bytecode, "", nil, 10); err != nil {
			t.Errorf("first execution failed: %v", err)
		}
	}()

	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, _, err := vm.ExecuteContext(ctx, bytecode, "", nil, 10)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected immediate error: %v", err)
		}
	case <-time.After(10 * time.Millisecond):
	}

	close(release)
	if err := <-done; err != nil {
		t.Fatalf("second execution failed: %v", err)
	}
	wg.Wait()
}

func TestSimpleVMStopCancelsQueuedExecutions(t *testing.T) {
	vm := NewSimpleVM(VMSuperLight)
	slowOpcode := uint32(0x0A0B0C)
	started := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	vm.RegisterOpcode(slowOpcode, func(in []byte) ([]byte, error) {
		once.Do(func() { close(started) })
		<-release
		return in, nil
	})
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	bytecode := []byte{0x0A, 0x0B, 0x0C}
	firstDone := make(chan error, 1)
	go func() {
		_, _, err := vm.Execute(bytecode, "", nil, 10)
		firstDone <- err
	}()

	<-started

	secondResult := make(chan error, 1)
	go func() {
		_, _, err := vm.Execute(bytecode, "", nil, 10)
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
		t.Fatalf("timed out waiting for queued execution cancellation")
	}

	close(release)
	if err := <-firstDone; err != nil {
		t.Fatalf("first execution failed: %v", err)
	}
	if err := <-stopErr; err != nil {
		t.Fatalf("stop failed: %v", err)
	}
}

func TestSimpleVMObserverReceivesStats(t *testing.T) {
	vm := NewSimpleVM()
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer vm.Stop()

	bytecode := []byte{0, 0, 0}
	observed := make(chan ExecutionStats, 1)
	vm.SetObserver(func(stats ExecutionStats) {
		observed <- stats
	})

	if _, _, err := vm.Execute(bytecode, "method", []byte{0xAA}, 10); err != nil {
		t.Fatalf("execute: %v", err)
	}

	select {
	case stats := <-observed:
		if stats.Method != "method" {
			t.Fatalf("unexpected method: %s", stats.Method)
		}
		if stats.GasUsed == 0 {
			t.Fatalf("expected gas usage")
		}
		if stats.Opcodes != 1 {
			t.Fatalf("expected 1 opcode, got %d", stats.Opcodes)
		}
	case <-time.After(time.Second):
		t.Fatalf("observer not invoked")
	}
}
