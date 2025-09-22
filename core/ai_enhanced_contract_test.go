package core

import (
	"testing"

	synnergy "synnergy"
)

func TestAIContractRegistry(t *testing.T) {
	vm := NewSimpleVM()
	_ = vm.Start()
	ledger := NewLedger()
	ledger.Credit("owner", 1_000_000)
	base := NewContractRegistry(vm, ledger)
	aiReg := NewAIContractRegistry(base)
	deployGas := synnergy.GasCost("DeployAIContract")
	addr, err := aiReg.DeployAIContract([]byte{0x01}, "abcd1234", "", deployGas, "owner")
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	if h, ok := aiReg.ModelHash(addr); !ok || h != "abcd1234" {
		t.Fatalf("model hash mismatch")
	}
	invokeGas := synnergy.GasCost("InvokeAIContract")
	out, _, err := aiReg.InvokeAIContract(addr, []byte("input"), invokeGas)
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	if string(out) != "input" {
		t.Fatalf("unexpected output")
	}
}
