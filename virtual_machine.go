package synnergy

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

// NewSimpleVM creates a new stopped virtual machine instance. An optional
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

	handlers := map[uint32]opcodeHandler{
		0x000000: func(b []byte) ([]byte, error) { // NOP/echo
			out := make([]byte, len(b))
			copy(out, b)
			return out, nil
		},
	}
	defaultH := handlers[0x000000]
	for _, op := range SNVMOpcodes {
		handlers[op.Code] = defaultH
	}

	vm := &SimpleVM{
		mode:     mode,
		limiter:  make(chan struct{}, capacity),
		handlers: handlers,
		defaultH: defaultH,
	}
	return vm
}

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

// Mode returns the configured resource profile for the VM.
func (vm *SimpleVM) Mode() VMMode {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.mode
}

// ModeString returns a human readable representation of the VM mode.
func (vm *SimpleVM) ModeString() string {
	switch vm.Mode() {
	case VMHeavy:
		return "heavy"
	case VMSuperLight:
		return "super-light"
	default:
		return "light"
	}
}

// RegisterOpcode associates an opcode with a handler. Nil handlers
// fallback to the VM's default no-op implementation.
func (vm *SimpleVM) RegisterOpcode(code uint32, h opcodeHandler) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if h == nil {
		h = vm.defaultH
	}
	vm.handlers[code] = h
}

// Execute interprets the provided bytecode as a sequence of 24-bit opcodes and
// dispatches to the registered handlers. Unknown opcodes fall back to a default
// echo handler so that tests remain deterministic. Gas is consumed per opcode.
func (vm *SimpleVM) ExecuteContext(ctx context.Context, wasm []byte, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !vm.Status() {
		return nil, 0, errors.New("vm not running")
	}

	// Bottleneck management: limit concurrent executions according to VM
	// profile. Respect context cancellation while waiting for a slot.
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

	opCount := uint64((len(wasm) + 2) / 3) // number of 24-bit opcodes
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
		handler, ok := vm.handlers[opcode]
		if !ok {
			handler = vm.defaultH
		}
		var err error
		out, err = handler(out)
		if err != nil {
			return nil, gasUsed, fmt.Errorf("opcode 0x%06x failed: %w", opcode, err)
		}
	}

	// simulate execution delay to keep behaviour deterministic
	select {
	case <-time.After(time.Millisecond):
	case <-ctx.Done():
		return nil, gasUsed, ctx.Err()
	}

	return out, gasUsed, nil
}

// Execute interprets the provided bytecode using a background context.
func (vm *SimpleVM) Execute(wasm []byte, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {
	return vm.ExecuteContext(context.Background(), wasm, method, args, gasLimit)
}
