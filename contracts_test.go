package synnergy_test

import (
	"testing"

	synnergy "synnergy"
	"synnergy/adapters/coreledger"
	"synnergy/core"
)

func TestContractRegistry(t *testing.T) {
	vm := synnergy.NewSimpleVM()
	if err := vm.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	ledger := core.NewLedger()
	ledger.Credit("owner", 100)
	reg := synnergy.NewContractRegistry(vm, coreledger.Wrap(ledger))
	wasm := []byte{0x00, 0x61}
	addr, err := reg.Deploy(wasm, "", 10, "owner")
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	if len(reg.List()) != 1 {
		t.Fatalf("expected 1 contract")
	}
	out, _, err := reg.Invoke(addr, "echo", []byte("hi"), 10)
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	if string(out) != "hi" {
		t.Fatalf("unexpected output: %q", out)
	}
}
