package core

import (
	"errors"
	"sync"
	"time"
)

// SimpleVM is a lightweight execution engine used by the contract registry. It
// mimics gas accounting and method dispatch but does not attempt to execute real
// WASM bytecode. The type satisfies the VirtualMachine interface.
type SimpleVM struct {
	mu      sync.RWMutex
	running bool
}

// NewSimpleVM creates a new stopped virtual machine instance.
func NewSimpleVM() *SimpleVM { return &SimpleVM{} }

// Start marks the VM as running. It is safe to call multiple times.
func (vm *SimpleVM) Start() error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.running {
		return nil
	}
	vm.running = true
	return nil
}

// Stop halts the VM instance.
func (vm *SimpleVM) Stop() error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if !vm.running {
		return nil
	}
	vm.running = false
	return nil
}

// Status reports whether the VM is running.
func (vm *SimpleVM) Status() bool {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.running
}

// Execute pretends to execute WASM bytecode. It returns the args unchanged as
// output and consumes gas proportionally to the input size. This behaviour keeps
// tests deterministic while providing a realistic API surface.
func (vm *SimpleVM) Execute(wasm []byte, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {
	if !vm.Status() {
		return nil, 0, errors.New("vm not running")
	}
	// naive gas model: 1 unit per byte of args with a minimum of 1
	gasUsed := uint64(len(args))
	if gasUsed == 0 {
		gasUsed = 1
	}
	if gasUsed > gasLimit {
		return nil, gasLimit, errors.New("gas limit exceeded")
	}
	// simulate execution delay
	time.Sleep(time.Millisecond)
	out := make([]byte, len(args))
	copy(out, args)
	return out, gasUsed, nil
}
