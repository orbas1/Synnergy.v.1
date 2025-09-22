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

// GasResolverFunc resolves the opcode name and gas cost for a numeric opcode.
// It allows the VM to integrate with the global gas catalogue while still
// supporting custom execution environments in tests.
type GasResolverFunc func(code uint32) (string, uint64)

// SimpleVM is a lightweight execution engine used by the contract registry. It
// interprets 24-bit opcodes from bytecode and dispatches them to handlers. The
// VM includes simple bottleneck management through a concurrency limiter and
// satisfies the VirtualMachine interface.
type SimpleVM struct {
	mu        sync.RWMutex
	running   bool
	mode      VMMode
	limiter   chan struct{}
	handlers  map[uint32]opcodeHandler
	defaultH  opcodeHandler
	gas       GasResolverFunc
	wg        sync.WaitGroup
	observer  ExecutionObserver
	lifecycle context.Context
	cancel    context.CancelFunc
}

// ExecutionStats capture metadata about a completed execution. They can be fed
// to external monitoring systems to provide visibility into VM behaviour in
// production deployments.
type ExecutionStats struct {
	Method         string
	GasUsed        uint64
	Opcodes        int
	Duration       time.Duration
	LastOpcode     uint32
	Mode           VMMode
	ContextDoneErr error
	PaddedOpcodes  int
}

// ExecutionObserver receives execution statistics after each successful run.
type ExecutionObserver func(ExecutionStats)

var (
	errVMNotRunning     = errors.New("vm not running")
	errVMBytecodeNeeded = errors.New("bytecode required")
	errVMGasLimit       = errors.New("gas limit exceeded")
)

var (
	opcodeOnce sync.Once
	opcodeIdx  map[uint32]string
)

func opcodeIndex() map[uint32]string {
	opcodeOnce.Do(func() {
		opcodeIdx = make(map[uint32]string, len(SNVMOpcodes)+len(ContractOpcodes))
		for _, op := range SNVMOpcodes {
			opcodeIdx[op.Code] = op.Name
		}
		for _, op := range ContractOpcodes {
			if _, exists := opcodeIdx[op.Code]; !exists {
				opcodeIdx[op.Code] = op.Name
			}
		}
	})
	return opcodeIdx
}

func defaultGasResolver(code uint32) (string, uint64) {
	if name, ok := opcodeIndex()[code]; ok {
		return name, GasCost(name)
	}
	return fmt.Sprintf("0x%06X", code), DefaultGasCost
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
		gas:      defaultGasResolver,
	}
	return vm
}

// SetGasResolver overrides the gas resolver used to price opcodes during
// execution. Providing nil restores the default resolver that integrates with
// the global gas catalogue.
func (vm *SimpleVM) SetGasResolver(resolver GasResolverFunc) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if resolver == nil {
		vm.gas = defaultGasResolver
		return
	}
	vm.gas = resolver
}

// SetObserver registers a callback to receive execution statistics whenever the
// VM completes a program. Passing nil clears the observer.
func (vm *SimpleVM) SetObserver(observer ExecutionObserver) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.observer = observer
}

// Start marks the VM as running. It is safe to call multiple times.
func (vm *SimpleVM) Start() error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.running {
		return nil
	}
	vm.lifecycle, vm.cancel = context.WithCancel(context.Background())
	vm.running = true
	return nil
}

// Stop halts the VM instance.
func (vm *SimpleVM) Stop() error {
	vm.mu.Lock()
	if !vm.running {
		vm.mu.Unlock()
		return nil
	}
	cancel := vm.cancel
	vm.cancel = nil
	vm.lifecycle = nil
	vm.running = false
	vm.observer = nil
	vm.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	vm.wg.Wait()
	return nil
}

// Status reports whether the VM is running.
func (vm *SimpleVM) Status() bool {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.running
}

// RegisterOpcode associates an opcode with a handler. Nil handlers
// fallback to the VM's default no-op implementation.
func (vm *SimpleVM) RegisterOpcode(code uint32, h opcodeHandler) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if h == nil {
		h = vm.defaultH
	}
	handlers := make(map[uint32]opcodeHandler, len(vm.handlers)+1)
	for k, v := range vm.handlers {
		handlers[k] = v
	}
	handlers[code] = h
	vm.handlers = handlers
}

// Execute interprets the provided bytecode as a sequence of 24-bit opcodes and
// dispatches to the registered handlers. Unknown opcodes fall back to a default
// echo handler so that tests remain deterministic. Gas is consumed per opcode.
func (vm *SimpleVM) ExecuteContext(ctx context.Context, wasm []byte, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	vm.mu.RLock()
	running := vm.running
	lifecycle := vm.lifecycle
	vm.mu.RUnlock()
	if !running {
		return nil, 0, errVMNotRunning
	}

	release, err := vm.acquireSlot(ctx, lifecycle)
	if err != nil {
		return nil, 0, err
	}
	defer release()

	vm.mu.RLock()
	if !vm.running {
		vm.mu.RUnlock()
		return nil, 0, errVMNotRunning
	}
	resolver := vm.gas
	handlers := vm.handlers
	defaultH := vm.defaultH
	observer := vm.observer
	mode := vm.mode
	vm.mu.RUnlock()
	if resolver == nil {
		resolver = defaultGasResolver
	}

	if len(wasm) == 0 {
		return nil, 0, errVMBytecodeNeeded
	}
	if gasLimit == 0 {
		return nil, 0, errVMGasLimit
	}

	out := make([]byte, len(args))
	copy(out, args)
	vm.wg.Add(1)
	defer vm.wg.Done()

	start := time.Now()
	var lastOpcode uint32
	var gasUsed uint64
	var opcodeCount int
	var paddedBytes int
	for i := 0; i < len(wasm); i += 3 {
		if err := ctx.Err(); err != nil {
			return nil, gasUsed, err
		}
		b0 := wasm[i]
		var b1, b2 byte
		if i+1 < len(wasm) {
			b1 = wasm[i+1]
		} else {
			paddedBytes++
		}
		if i+2 < len(wasm) {
			b2 = wasm[i+2]
		} else {
			paddedBytes++
		}
		opcode := uint32(b0)<<16 | uint32(b1)<<8 | uint32(b2)
		lastOpcode = opcode
		_, cost := resolver(opcode)
		if cost == 0 {
			cost = DefaultGasCost
		}
		if gasUsed+cost > gasLimit {
			return nil, gasUsed, errVMGasLimit
		}
		handler, ok := handlers[opcode]
		if !ok {
			handler = defaultH
		}
		var err error
		out, err = handler(out)
		if err != nil {
			return nil, gasUsed, fmt.Errorf("opcode 0x%06x failed: %w", opcode, err)
		}
		gasUsed += cost
		opcodeCount++
	}

	duration := time.Since(start)

	if err := ctx.Err(); err != nil {
		return nil, gasUsed, err
	}

	if observer != nil {
		observer(ExecutionStats{
			Method:         method,
			GasUsed:        gasUsed,
			Opcodes:        opcodeCount,
			Duration:       duration,
			LastOpcode:     lastOpcode,
			Mode:           mode,
			ContextDoneErr: ctx.Err(),
			PaddedOpcodes:  paddedBytes,
		})
	}

	return out, gasUsed, nil
}

// Execute interprets the provided bytecode using a background context.
func (vm *SimpleVM) Execute(wasm []byte, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {
	return vm.ExecuteContext(context.Background(), wasm, method, args, gasLimit)
}

func (vm *SimpleVM) acquireSlot(ctx context.Context, lifecycle context.Context) (func(), error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var stop <-chan struct{}
	if lifecycle != nil {
		stop = lifecycle.Done()
	}
	select {
	case vm.limiter <- struct{}{}:
		return func() { <-vm.limiter }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-stop:
		return nil, errVMNotRunning
	}
}
