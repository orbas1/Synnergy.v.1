package core

import "testing"

func TestAIContractRegistry(t *testing.T) {
	vm := NewSimpleVM()
	_ = vm.Start()
	base := NewContractRegistry(vm)
	aiReg := NewAIContractRegistry(base)
	addr, err := aiReg.DeployAIContract([]byte{0x01}, "modelhash", "", 5, "owner")
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	if h, ok := aiReg.ModelHash(addr); !ok || h != "modelhash" {
		t.Fatalf("model hash mismatch")
	}
	out, _, err := aiReg.InvokeAIContract(addr, []byte("input"), 5)
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	if string(out) != "input" {
		t.Fatalf("unexpected output")
	}
}
