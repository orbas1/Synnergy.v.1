package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
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

// String returns a descriptive label for the VM profile.
func (m VMMode) String() string {
	switch m {
	case VMHeavy:
		return "heavy"
	case VMLight:
		return "light"
	case VMSuperLight:
		return "super_light"
	default:
		return "unknown"
	}
}

// opcodeHandler defines the function signature for executing a single
// opcode. Input bytes are transformed and returned as new output.
type opcodeHandler func([]byte) ([]byte, error)

// SimpleVM is a lightweight execution engine used by the contract registry. It
// interprets 24-bit opcodes from bytecode and dispatches them to handlers. The
// VM includes simple bottleneck management through a concurrency limiter and
// satisfies the VirtualMachine interface.
type SimpleVM struct {
	mu           sync.RWMutex
	running      bool
	mode         VMMode
	limiter      chan struct{}
	handlers     map[uint32]opcodeHandler
	defaultH     opcodeHandler
	callHandlers map[string]func() error
	hooks        []ExecutionHook
	metrics      vmMetrics
}

// ErrGasLimit is returned when execution exhausts its allotted gas.
var ErrGasLimit = errors.New("gas limit exceeded")

// ExecutionTrace captures a single opcode execution, allowing observers to
// plug into the VM for telemetry, auditing or debugging. Stage 75 introduces
// the trace model so both the CLI and the JavaScript console can surface
// real-time execution insights.
type ExecutionTrace struct {
	Opcode       uint32
	Name         string
	GasCost      uint64
	Duration     time.Duration
	RemainingGas uint64
	Err          error
}

// ExecutionHook is invoked after each opcode completes. Hooks should be fast
// and never panic; the VM guards against panics to preserve fault tolerance.
type ExecutionHook func(ExecutionTrace)

type vmMetrics struct {
	executions uint64
	failures   uint64
	gasUsed    uint64
	durationNS int64
}

// VMMetrics exposes aggregated statistics for observability dashboards and
// CLI diagnostics.
type VMMetrics struct {
	Executions    uint64
	Failures      uint64
	GasConsumed   uint64
	TotalDuration time.Duration
}

type executionContext struct {
	vm        *SimpleVM
	limit     uint64
	remaining uint64
	used      uint64
}

// Call implements the OpContext interface used by the global opcode dispatcher.
// For now it simply acts as a stub; higher-level integration can wire this to
// actual protocol functionality.
func (vm *SimpleVM) Call(name string) error {
	vm.mu.RLock()
	handler := vm.callHandlers[name]
	vm.mu.RUnlock()
	if handler == nil {
		return nil
	}
	return handler()
}

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
		mode:         mode,
		limiter:      make(chan struct{}, capacity),
		callHandlers: make(map[string]func() error),
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

// RegisterCallHandler wires a named call target to a callback. This allows the
// dispatcher to invoke high-level protocol handlers through the VM façade
// without coupling tests to specific implementations.
func (vm *SimpleVM) RegisterCallHandler(name string, fn func() error) {
	vm.mu.Lock()
	if vm.callHandlers == nil {
		vm.callHandlers = make(map[string]func() error)
	}
	if fn == nil {
		delete(vm.callHandlers, name)
	} else {
		vm.callHandlers[name] = fn
	}
	vm.mu.Unlock()
}

// RegisterHook attaches an execution hook that receives opcode traces. Hooks
// are invoked sequentially after each opcode completes.
func (vm *SimpleVM) RegisterHook(h ExecutionHook) {
	if h == nil {
		return
	}
	vm.mu.Lock()
	vm.hooks = append(vm.hooks, h)
	vm.mu.Unlock()
}

// ResetHooks removes all registered execution hooks. Intended for tests to
// avoid cross-test interference.
func (vm *SimpleVM) ResetHooks() {
	vm.mu.Lock()
	vm.hooks = nil
	vm.mu.Unlock()
}

// ResetMetrics zeroes the VM statistics. Primarily used in tests and
// diagnostic tooling to capture deltas.
func (vm *SimpleVM) ResetMetrics() {
	atomic.StoreUint64(&vm.metrics.executions, 0)
	atomic.StoreUint64(&vm.metrics.failures, 0)
	atomic.StoreUint64(&vm.metrics.gasUsed, 0)
	atomic.StoreInt64(&vm.metrics.durationNS, 0)
}

// Metrics returns a snapshot of the VM statistics.
func (vm *SimpleVM) Metrics() VMMetrics {
	return VMMetrics{
		Executions:    atomic.LoadUint64(&vm.metrics.executions),
		Failures:      atomic.LoadUint64(&vm.metrics.failures),
		GasConsumed:   atomic.LoadUint64(&vm.metrics.gasUsed),
		TotalDuration: time.Duration(atomic.LoadInt64(&vm.metrics.durationNS)),
	}
}

func (vm *SimpleVM) snapshotHooks() []ExecutionHook {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	if len(vm.hooks) == 0 {
		return nil
	}
	hooks := make([]ExecutionHook, len(vm.hooks))
	copy(hooks, vm.hooks)
	return hooks
}

func (vm *SimpleVM) emitTrace(hooks []ExecutionHook, trace ExecutionTrace) {
	for _, h := range hooks {
		func() {
			defer func() {
				_ = recover()
			}()
			h(trace)
		}()
	}
}

func (vm *SimpleVM) recordMetrics(gas uint64, duration time.Duration, err error) {
	atomic.AddUint64(&vm.metrics.executions, 1)
	atomic.AddUint64(&vm.metrics.gasUsed, gas)
	atomic.AddInt64(&vm.metrics.durationNS, duration.Nanoseconds())
	if err != nil {
		atomic.AddUint64(&vm.metrics.failures, 1)
	}
}

func newExecutionContext(vm *SimpleVM, gasLimit uint64) *executionContext {
	return &executionContext{vm: vm, limit: gasLimit, remaining: gasLimit}
}

func (ec *executionContext) Call(name string) error {
	return ec.vm.Call(name)
}

func (ec *executionContext) Gas(amount uint64) error {
	if amount == 0 {
		return nil
	}
	if amount > ec.remaining {
		ec.used += amount
		if ec.used > ec.limit {
			ec.used = ec.limit
		}
		ec.remaining = 0
		return fmt.Errorf("%w: required %d remaining %d", ErrGasLimit, amount, 0)
	}
	ec.remaining -= amount
	ec.used += amount
	return nil
}

func (ec *executionContext) Remaining() uint64 { return ec.remaining }

func (ec *executionContext) Used() uint64 { return ec.used }

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

// Mode reports the configured VM profile.
func (vm *SimpleVM) Mode() VMMode {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.mode
}

// ExecuteContext interprets bytecode with context cancellation and gas limits.
func (vm *SimpleVM) ExecuteContext(ctx context.Context, wasm []byte, method string, args []byte, gasLimit uint64) (out []byte, gasUsed uint64, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !vm.Status() {
		return nil, 0, errors.New("vm not running")
	}

	release, err := vm.acquireSlot(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer release()

	if len(wasm) == 0 {
		return nil, 0, errors.New("bytecode required")
	}

	hooks := vm.snapshotHooks()
	exec := newExecutionContext(vm, gasLimit)
	start := time.Now()
	out = args

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("vm panic: %v", r)
		}
		gasUsed = exec.Used()
		vm.recordMetrics(gasUsed, time.Since(start), err)
	}()

	opcodeNames := Opcodes()
	for i := 0; i < len(wasm); i += 3 {
		if ctx.Err() != nil {
			err = ctx.Err()
			return nil, exec.Used(), err
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

		cost := vm.gasCostForOpcode(opcode)
		name := vm.nameForOpcode(opcode, opcodeNames)
		stepStart := time.Now()

		if handler, ok := vm.handlers[opcode]; ok {
			if err = exec.Gas(cost); err != nil {
				vm.emitTrace(hooks, ExecutionTrace{Opcode: opcode, Name: name, GasCost: cost, Duration: time.Since(stepStart), RemainingGas: exec.Remaining(), Err: err})
				return nil, exec.Used(), err
			}
			out, err = handler(out)
			vm.emitTrace(hooks, ExecutionTrace{Opcode: opcode, Name: name, GasCost: cost, Duration: time.Since(stepStart), RemainingGas: exec.Remaining(), Err: err})
			if err != nil {
				err = fmt.Errorf("opcode 0x%06x failed: %w", opcode, err)
				return nil, exec.Used(), err
			}
			continue
		}

		if opcode == 0 {
			if err = exec.Gas(cost); err != nil {
				vm.emitTrace(hooks, ExecutionTrace{Opcode: opcode, Name: name, GasCost: cost, Duration: time.Since(stepStart), RemainingGas: exec.Remaining(), Err: err})
				return nil, exec.Used(), err
			}
			if vm.defaultH != nil {
				out, err = vm.defaultH(out)
				vm.emitTrace(hooks, ExecutionTrace{Opcode: opcode, Name: name, GasCost: cost, Duration: time.Since(stepStart), RemainingGas: exec.Remaining(), Err: err})
				if err != nil {
					err = fmt.Errorf("opcode 0x%06x failed: %w", opcode, err)
					return nil, exec.Used(), err
				}
			} else {
				vm.emitTrace(hooks, ExecutionTrace{Opcode: opcode, Name: name, GasCost: cost, Duration: time.Since(stepStart), RemainingGas: exec.Remaining(), Err: nil})
			}
			continue
		}

		err = Dispatch(exec, Opcode(opcode))
		if err != nil {
			vm.emitTrace(hooks, ExecutionTrace{Opcode: opcode, Name: name, GasCost: cost, Duration: time.Since(stepStart), RemainingGas: exec.Remaining(), Err: err})
			continue
		}
		vm.emitTrace(hooks, ExecutionTrace{Opcode: opcode, Name: name, GasCost: cost, Duration: time.Since(stepStart), RemainingGas: exec.Remaining(), Err: nil})
	}

	if err = vm.awaitDeterministicDelay(ctx); err != nil {
		return nil, exec.Used(), err
	}
	return out, exec.Used(), nil
}

func (vm *SimpleVM) Execute(wasm []byte, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {
	return vm.ExecuteContext(context.Background(), wasm, method, args, gasLimit)
}

func (vm *SimpleVM) acquireSlot(ctx context.Context) (func(), error) {
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case vm.limiter <- struct{}{}:
		return func() { <-vm.limiter }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, errors.New("vm busy")
	}
}

func (vm *SimpleVM) awaitDeterministicDelay(ctx context.Context) error {
	select {
	case <-time.After(time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (vm *SimpleVM) gasCostForOpcode(opcode uint32) uint64 {
	cost := GasCost(Opcode(opcode))
	if cost == 0 {
		return DefaultGasCost
	}
	return cost
}

func (vm *SimpleVM) nameForOpcode(opcode uint32, names map[Opcode]string) string {
	if names != nil {
		if name, ok := names[Opcode(opcode)]; ok {
			return name
		}
	}
	if opcode == 0 {
		return "NOP"
	}
	return fmt.Sprintf("0x%06X", opcode)
}
