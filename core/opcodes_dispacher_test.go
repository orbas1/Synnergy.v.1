package core

import "testing"

type mockCtx struct {
	called []string
	gas    uint64
}

func (m *mockCtx) Call(name string) error {
	m.called = append(m.called, name)
	return nil
}

func (m *mockCtx) Gas(g uint64) error {
	m.gas += g
	return nil
}

// TestBindVM ensures that opcode handlers are registered with the VM and
// delegate through the dispatcher using the provided context.
func TestBindVM(t *testing.T) {
	vm := NewSimpleVM()
	ctx := &mockCtx{}
	BindVM(vm, func() OpContext { return ctx })

	b, err := ToBytecode("InitContracts")
	if err != nil {
		t.Fatalf("opcode lookup failed: %v", err)
	}
	op, err := ParseOpcode(b)
	if err != nil {
		t.Fatalf("opcode parse failed: %v", err)
	}

	h, ok := vm.handlers[uint32(op)]
	if !ok {
		t.Fatalf("handler for opcode %#x not registered", op)
	}

	if _, err := h(nil); err != nil {
		t.Fatalf("handler execution failed: %v", err)
	}

	if len(ctx.called) == 0 {
		t.Fatalf("expected context Call to be invoked")
	}
}
