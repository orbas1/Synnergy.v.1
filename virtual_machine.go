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

type opcodeBinding struct {
	name    string
	handler opcodeHandler
}

// SimpleVM is a lightweight execution engine used by the contract registry. It
// interprets 24-bit opcodes from bytecode and dispatches them to handlers. The
// VM includes simple bottleneck management through a concurrency limiter and
// satisfies the VirtualMachine interface.
type SimpleVM struct {
	mu             sync.RWMutex
	running        bool
	mode           VMMode
	limiter        chan struct{}
	handlers       map[uint32]opcodeBinding
	defaultBinding opcodeBinding

	gasMu          sync.RWMutex
	gasTable       GasTable
	gasOverrides   map[string]uint64
	gasVersion     uint64
	gasLoadedAt    time.Time
	gasUpdates     <-chan GasSnapshot
	cancelGasWatch func()
	gasWatcherDone chan struct{}
	closeOnce      sync.Once
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
	for _, op := range SNVMOpcodes {
		if existing, ok := handlers[op.Code]; ok {
			if existing.name == "" {
				existing.name = op.Name
			}
			if existing.handler == nil {
				existing.handler = defaultBinding.handler
			}
			handlers[op.Code] = existing
			continue
		}
		handlers[op.Code] = opcodeBinding{name: op.Name, handler: defaultBinding.handler}
	}

	vm := &SimpleVM{
		mode:           mode,
		limiter:        make(chan struct{}, capacity),
		handlers:       handlers,
		defaultBinding: defaultBinding,
		gasOverrides:   make(map[string]uint64),
		gasWatcherDone: make(chan struct{}),
	}

	snapshot := SnapshotGasTable()
	vm.gasTable = snapshot.Table
	for name, cost := range snapshot.Table {
		if cost != 0 {
			vm.gasOverrides[name] = cost
		}
	}
	vm.gasVersion = snapshot.Version
	vm.gasLoadedAt = snapshot.LoadedAt
	updates, cancel := SubscribeGasTable(4)
	vm.gasUpdates = updates
	vm.cancelGasWatch = cancel
	go vm.watchGasUpdates()

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

// RegisterOpcode associates an opcode with a handler. Nil handlers
// fallback to the VM's default no-op implementation.
func (vm *SimpleVM) RegisterOpcode(code uint32, h opcodeHandler) {
	vm.RegisterOpcodeNamed(code, "", h)
}

// RegisterOpcodeNamed allows callers to bind a specific identifier to the
// opcode. When no name is provided the existing binding is retained or a
// deterministic placeholder is generated so gas accounting remains accurate.
func (vm *SimpleVM) RegisterOpcodeNamed(code uint32, name string, h opcodeHandler) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if h == nil {
		h = vm.defaultBinding.handler
	}
	binding, ok := vm.handlers[code]
	if !ok {
		binding = opcodeBinding{}
	}
	if name != "" {
		binding.name = name
	} else if binding.name == "" {
		binding.name = fmt.Sprintf("opcode_0x%06x", code)
	}
	binding.handler = h
	vm.handlers[code] = binding
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
			cost = DefaultGasCost
		}
		if gasLimit < cost || gasUsed > gasLimit-cost {
			return nil, gasLimit, errors.New("gas limit exceeded")
		}
		gasUsed += cost
		var err error
		out, err = binding.handler(out)
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

func (vm *SimpleVM) bindingFor(code uint32) opcodeBinding {
	vm.mu.RLock()
	binding, ok := vm.handlers[code]
	vm.mu.RUnlock()
	if !ok {
		return vm.defaultBinding
	}
	if binding.handler == nil {
		binding.handler = vm.defaultBinding.handler
	}
	if binding.name == "" {
		binding.name = fmt.Sprintf("opcode_0x%06x", code)
	}
	return binding
}

func (vm *SimpleVM) gasFor(name string) uint64 {
	if name == "" {
		return DefaultGasCost
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
	cost := GasCost(name)
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
func (vm *SimpleVM) GasSnapshot() GasSnapshot {
	vm.gasMu.RLock()
	defer vm.gasMu.RUnlock()
	return GasSnapshot{Table: cloneGasTable(vm.gasTable), Version: vm.gasVersion, LoadedAt: vm.gasLoadedAt}
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
