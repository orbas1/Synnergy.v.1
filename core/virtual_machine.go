package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// VMMode represents the resource profile of the virtual machine.
// Heavy instances allow more concurrent executions while super light
// instances are tailored for constrained environments like mobile
// devices.
type VMMode int

const (
	VMHeavy VMMode = iota
	VMLight
	VMSuperLight
)

// opcodeHandler defines the function signature for executing a single
// opcode. Input bytes are transformed and returned as new output.
type opcodeHandler func([]byte) ([]byte, error)

// SimpleVM is a lightweight execution engine used by the contract registry. It
// interprets 24-bit opcodes from bytecode and dispatches them to handlers. The
// VM includes simple bottleneck management through a concurrency limiter and
// satisfies the VirtualMachine interface.
type SimpleVM struct {
	mu       sync.RWMutex
	running  bool
	mode     VMMode
	limiter  chan struct{}
	handlers map[uint32]opcodeHandler
	defaultH opcodeHandler
}

// Call implements the OpContext interface used by the global opcode dispatcher.
// For now it simply acts as a stub; higher-level integration can wire this to
// actual protocol functionality.
func (vm *SimpleVM) Call(string) error { return nil }

// Gas satisfies the OpContext interface. The lightweight VM does not meter gas
// beyond counting opcodes, so this is a no-op placeholder.
func (vm *SimpleVM) Gas(uint64) error { return nil }

// NewSimpleVM creates a new stopped virtual machine instance.  An optional
// mode can be supplied to configure resource limits; by default a light VM is
// created.
func NewSimpleVM(modes ...VMMode) *SimpleVM {
	mode := VMLight
	if len(modes) > 0 {
		mode = modes[0]
	}

	// Concurrency capacities for different VM profiles. The heavy VM allows
	// more parallel executions while the super light VM processes a single
	// request at a time.
	capacity := 5
	switch mode {
	case VMHeavy:
		capacity = 10
	case VMSuperLight:
		capacity = 1
	}

	vm := &SimpleVM{
		mode:    mode,
		limiter: make(chan struct{}, capacity),
	}
	vm.handlers = map[uint32]opcodeHandler{
		0x000000: func(b []byte) ([]byte, error) { // NOP/echo
			out := make([]byte, len(b))
			copy(out, b)
			return out, nil
		},
	}
	vm.defaultH = vm.handlers[0x000000]
	return vm
}

// RegisterHandler allows callers to inject or override opcode handlers at
// runtime.  It is safe for concurrent use and enables higher level modules or
// tests to extend the VM without recompilation.
func (vm *SimpleVM) RegisterHandler(op uint32, h opcodeHandler) {
	vm.mu.Lock()
	if vm.handlers == nil {
		vm.handlers = make(map[uint32]opcodeHandler)
	}
	vm.handlers[op] = h
	vm.mu.Unlock()
}

// Concurrency returns the maximum number of executions the VM will run in
// parallel based on its current limiter capacity.  This is primarily exported
// for monitoring and tests.
func (vm *SimpleVM) Concurrency() int { return cap(vm.limiter) }

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

// ExecuteContext interprets bytecode with context cancellation and gas limits.
func (vm *SimpleVM) ExecuteContext(ctx context.Context, wasm []byte, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !vm.Status() {
		return nil, 0, errors.New("vm not running")
	}

	select {
	case vm.limiter <- struct{}{}:
		defer func() { <-vm.limiter }()
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
		return nil, 0, errors.New("vm busy")
	}

	if len(wasm) == 0 {
		return nil, 0, errors.New("bytecode required")
	}

	opCount := (uint64(len(wasm)) + 2) / 3
	gasUsed := opCount
	if gasUsed == 0 {
		gasUsed = 1
	}
	if gasUsed > gasLimit {
		return nil, gasLimit, errors.New("gas limit exceeded")
	}

	out := args
	for i := 0; i < len(wasm); i += 3 {
		if err := ctx.Err(); err != nil {
			return nil, gasUsed, err
		}
		b0 := wasm[i]
		var b1, b2 byte
		if i+1 < len(wasm) {
			b1 = wasm[i+1]
		}
		if i+2 < len(wasm) {
			b2 = wasm[i+2]
		}
		opcode := uint32(b0)<<16 | uint32(b1)<<8 | uint32(b2)

		if handler, ok := vm.handlers[opcode]; ok {
			var err error
			out, err = handler(out)
			if err != nil {
				return nil, gasUsed, fmt.Errorf("opcode 0x%06x failed: %w", opcode, err)
			}
			continue
		}
		if opcode != 0 {
			if err := Dispatch(vm, Opcode(opcode)); err != nil {
				continue
			}
		}
	}

	select {
	case <-time.After(time.Millisecond):
	case <-ctx.Done():
		return nil, gasUsed, ctx.Err()
	}

	return out, gasUsed, nil
}

func (vm *SimpleVM) Execute(wasm []byte, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {
	return vm.ExecuteContext(context.Background(), wasm, method, args, gasLimit)
}
