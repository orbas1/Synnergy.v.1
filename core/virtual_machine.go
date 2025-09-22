package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	synn "synnergy"
)

// VMMode represents the resource profile of the virtual machine.
// Heavy instances allow more concurrent executions while super light
// instances are tailored for constrained environments like mobile devices.
type VMMode int

const (
	VMHeavy VMMode = iota
	VMLight
	VMSuperLight
)

type opcodeHandler func([]byte) ([]byte, error)

type opcodeBinding struct {
	name     string
	handler  opcodeHandler
	dispatch bool
}

// SimpleVM is the execution engine used by the contract registry and CLI
// helpers. It mirrors the root VM implementation but wires in the opcode
// dispatcher so unknown codes are resolved through the generated catalogue.
type SimpleVM struct {
	mu             sync.RWMutex
	running        bool
	mode           VMMode
	limiter        chan struct{}
	handlers       map[uint32]opcodeBinding
	defaultBinding opcodeBinding

	gasMu          sync.RWMutex
	gasTable       synn.GasTable
	gasOverrides   map[string]uint64
	gasVersion     uint64
	gasLoadedAt    time.Time
	gasUpdates     <-chan synn.GasSnapshot
	cancelGasWatch func()
	gasWatcherDone chan struct{}
	closeOnce      sync.Once
}

// Call implements the OpContext interface used by the opcode dispatcher. The
// lightweight VM does not expose additional host functionality yet.
func (vm *SimpleVM) Call(string) error { return nil }

// Gas satisfies the OpContext interface. The lightweight VM meters gas locally
// rather than delegating to the dispatcher, so this is a no-op.
func (vm *SimpleVM) Gas(uint64) error { return nil }

// Concurrency returns the number of parallel executions permitted by the
// limiter. It provides visibility for monitoring and tests.
func (vm *SimpleVM) Concurrency() int { return cap(vm.limiter) }

// NewSimpleVM creates a new stopped virtual machine instance. An optional mode
// can be supplied to configure resource limits; by default a light VM is created.
func NewSimpleVM(modes ...VMMode) *SimpleVM {
	mode := VMLight
	if len(modes) > 0 {
		mode = modes[0]
	}

	capacity := 5
	switch mode {
	case VMHeavy:
		capacity = 10
	case VMSuperLight:
		capacity = 1
	}

	defaultBinding := opcodeBinding{
		name: "NOP",
		handler: func(b []byte) ([]byte, error) {
			out := make([]byte, len(b))
			copy(out, b)
			return out, nil
		},
	}

	handlers := map[uint32]opcodeBinding{
		0x000000: defaultBinding,
	}

	vm := &SimpleVM{
		mode:           mode,
		limiter:        make(chan struct{}, capacity),
		handlers:       handlers,
		defaultBinding: defaultBinding,
		gasOverrides:   make(map[string]uint64),
		gasWatcherDone: make(chan struct{}),
	}

	snapshot := synn.SnapshotGasTable()
	vm.gasTable = snapshot.Table
	for name, cost := range snapshot.Table {
		if cost != 0 {
			vm.gasOverrides[name] = cost
		}
	}
	vm.gasVersion = snapshot.Version
	vm.gasLoadedAt = snapshot.LoadedAt
	updates, cancel := synn.SubscribeGasTable(4)
	vm.gasUpdates = updates
	vm.cancelGasWatch = cancel
	go vm.watchGasUpdates()

	return vm
}

// RegisterHandler allows callers to inject or override opcode handlers at
// runtime. It is safe for concurrent use and enables higher level modules or
// tests to extend the VM without recompilation.
func (vm *SimpleVM) RegisterHandler(op uint32, h opcodeHandler) {
	vm.RegisterHandlerNamed(op, fmt.Sprintf("custom_0x%06x", op), h)
}

// RegisterHandlerNamed allows callers to inject opcode handlers with explicit gas labels.
func (vm *SimpleVM) RegisterHandlerNamed(op uint32, name string, h opcodeHandler) {
	vm.mu.Lock()
	if vm.handlers == nil {
		vm.handlers = make(map[uint32]opcodeBinding)
	}
	vm.handlers[op] = opcodeBinding{name: name, handler: h}
	vm.mu.Unlock()
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

// Close releases background watchers. It is safe to invoke multiple times and
// should be called when the VM will no longer be used.
func (vm *SimpleVM) Close() {
	vm.closeOnce.Do(func() {
		if vm.cancelGasWatch != nil {
			vm.cancelGasWatch()
		}
		if vm.gasWatcherDone != nil {
			<-vm.gasWatcherDone
		}
	})
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

	var gasUsed uint64
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
		binding := vm.bindingFor(opcode)
		cost := vm.gasFor(binding.name)
		if cost == 0 {
			cost = synn.DefaultGasCost
		}
		if gasLimit < cost || gasUsed > gasLimit-cost {
			return nil, gasLimit, errors.New("gas limit exceeded")
		}
		gasUsed += cost

		if binding.dispatch {
			if err := Dispatch(vm, Opcode(opcode)); err != nil {
				return nil, gasUsed, fmt.Errorf("dispatch 0x%06x failed: %w", opcode, err)
			}
			continue
		}

		var err error
		out, err = binding.handler(out)
		if err != nil {
			return nil, gasUsed, fmt.Errorf("opcode 0x%06x failed: %w", opcode, err)
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

func (vm *SimpleVM) bindingFor(code uint32) opcodeBinding {
	vm.mu.RLock()
	binding, ok := vm.handlers[code]
	vm.mu.RUnlock()
	if ok {
		if binding.handler == nil {
			binding.handler = vm.defaultBinding.handler
		}
		if binding.name == "" {
			binding.name = fmt.Sprintf("opcode_0x%06x", code)
		}
		return binding
	}
	name, dispatch := vm.lookupOpcodeName(code)
	if !dispatch {
		return opcodeBinding{name: name, handler: vm.defaultBinding.handler}
	}
	if name == "" {
		name = fmt.Sprintf("opcode_0x%06x", code)
	}
	return opcodeBinding{name: name, handler: vm.defaultBinding.handler, dispatch: true}
}

func (vm *SimpleVM) lookupOpcodeName(code uint32) (string, bool) {
	if op, ok := synn.SNVMOpcodeByCode(code); ok {
		return op.Name, true
	}
	if name, ok := LookupOpcodeName(code); ok {
		return name, true
	}
	return fmt.Sprintf("opcode_0x%06x", code), false
}

func (vm *SimpleVM) gasFor(name string) uint64 {
	if name == "" {
		return synn.DefaultGasCost
	}
	vm.gasMu.RLock()
	table := vm.gasTable
	override := vm.gasOverrides[name]
	vm.gasMu.RUnlock()
	if table != nil {
		if cost, ok := table[name]; ok {
			vm.recordGasOverride(name, cost)
			return cost
		}
	}
	if override != 0 {
		return override
	}
	cost := synn.GasCost(name)
	vm.recordGasOverride(name, cost)
	return cost
}

func (vm *SimpleVM) watchGasUpdates() {
	if vm.gasUpdates == nil {
		close(vm.gasWatcherDone)
		return
	}
	defer close(vm.gasWatcherDone)
	for snapshot := range vm.gasUpdates {
		vm.gasMu.Lock()
		vm.gasTable = snapshot.Table
		vm.gasVersion = snapshot.Version
		vm.gasLoadedAt = snapshot.LoadedAt
		if vm.gasOverrides == nil {
			vm.gasOverrides = make(map[string]uint64)
		}
		for name, cost := range snapshot.Table {
			if cost != 0 {
				vm.gasOverrides[name] = cost
			}
		}
		vm.gasMu.Unlock()
	}
}

// GasSnapshot exposes the most recent gas table cached by the VM. It allows
// external tooling and tests to validate that the VM is operating with the
// expected pricing data without reaching into package internals.
func (vm *SimpleVM) GasSnapshot() synn.GasSnapshot {
	vm.gasMu.RLock()
	defer vm.gasMu.RUnlock()
	table := make(synn.GasTable, len(vm.gasTable))
	for k, v := range vm.gasTable {
		table[k] = v
	}
	return synn.GasSnapshot{Table: table, Version: vm.gasVersion, LoadedAt: vm.gasLoadedAt}
}

func (vm *SimpleVM) recordGasOverride(name string, cost uint64) {
	if name == "" || cost == 0 {
		return
	}
	vm.gasMu.Lock()
	if vm.gasOverrides == nil {
		vm.gasOverrides = make(map[string]uint64)
	}
	vm.gasOverrides[name] = cost
	vm.gasMu.Unlock()
}

// LookupOpcodeName bridges the dispatcher catalogue to provide human readable
// names for opcodes handled externally.
func LookupOpcodeName(code uint32) (string, bool) {
	ops := Opcodes()
	if name, ok := ops[Opcode(code)]; ok {
		return name, true
	}
	return "", false
}
