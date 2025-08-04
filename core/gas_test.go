package core

import "testing"

func TestDefaultGasTable(t *testing.T) {
    g := DefaultGasTable()
    if g[OpNoop] != 1 {
        t.Fatalf("expected gas 1 for OpNoop")
    }
    if g[OpTransfer] != 10 {
        t.Fatalf("expected gas 10 for OpTransfer")
    }
    if g[OpAdd] != 0 {
        t.Fatalf("expected undefined ops to have zero gas cost")
    }
}

