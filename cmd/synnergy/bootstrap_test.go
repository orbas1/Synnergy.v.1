package main

import (
	"context"
	"testing"

	synn "synnergy"
)

func TestEnsureGasCatalogue(t *testing.T) {
	synn.ResetGasTable()
	if err := ensureGasCatalogue(); err != nil {
		t.Fatalf("ensureGasCatalogue: %v", err)
	}

	for _, op := range requiredGasOperations {
		if !synn.HasOpcode(op) {
			t.Fatalf("expected opcode %s to be registered", op)
		}
	}
}

func TestBootstrapRuntimeLifecycle(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rt, err := bootstrapRuntime(ctx)
	if err != nil {
		t.Fatalf("bootstrapRuntime: %v", err)
	}
	if rt == nil || rt.vm == nil {
		t.Fatal("runtime missing VM")
	}
	if !rt.vm.Status() {
		t.Fatal("VM should be running after bootstrap")
	}

	rt.Shutdown()
	if rt.vm.Status() {
		t.Fatal("VM should be stopped after shutdown")
	}
}
