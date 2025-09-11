package core

import "testing"

// mock context implementing OpContext for tests
type testCtx struct {
	called []string
	gas    uint64
}

func (c *testCtx) Call(name string) error {
	c.called = append(c.called, name)
	return nil
}
func (c *testCtx) Gas(n uint64) error {
	c.gas += n
	return nil
}

func TestLookupAndDispatch(t *testing.T) {
	cat := Catalogue()
	if len(cat) == 0 {
		t.Skip("no opcodes registered")
	}
	entry := cat[0]
	op, ok := Lookup(entry.Name)
	if !ok || op != entry.Op {
		t.Fatalf("lookup mismatch: %v vs %v", op, entry.Op)
	}
	ctx := &testCtx{}
	if err := Dispatch(ctx, op); err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if len(ctx.called) != 1 || ctx.called[0] != entry.Name {
		t.Fatalf("handler not invoked: %v", ctx.called)
	}
	if ctx.gas == 0 {
		t.Fatalf("gas not charged")
	}
}
