package core

// BindVM integrates the opcode dispatcher with a SimpleVM instance. It
// registers a handler for every known opcode so that smart contract bytecode
// executed by the VM can invoke high-level protocol functions. The provided
// ctxFn is called for each dispatched opcode and should return an OpContext
// configured for that invocation (e.g. with ledger references or gas meters).
// If ctxFn is nil, a no-op context is used which simply treats all operations
// as free and performs no side effects.
func BindVM(vm *SimpleVM, ctxFn func() OpContext) {
	if ctxFn == nil {
		ctxFn = func() OpContext { return noopContext{} }
	}

	vm.mu.Lock()
	defer vm.mu.Unlock()

	if vm.handlers == nil {
		vm.handlers = make(map[uint32]opcodeHandler, len(opcodeTable))
	}

	for op := range opcodeTable {
		opCopy := op
		vm.handlers[uint32(opCopy)] = func(in []byte) ([]byte, error) {
			if err := Dispatch(ctxFn(), opCopy); err != nil {
				return nil, err
			}
			return in, nil
		}
	}
}

// noopContext satisfies OpContext and is used when no context factory is
// supplied. All methods are no-ops, allowing opcode registration in tests or
// standalone tools without needing a full blockchain environment.
type noopContext struct{}

func (noopContext) Call(string) error { return nil }
func (noopContext) Gas(uint64) error  { return nil }
