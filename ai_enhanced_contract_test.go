package synnergy

import "testing"

func TestAIContractRegistry(t *testing.T) {
	vm := NewSimpleVM()
	_ = vm.Start()
	reg := NewContractRegistry(vm)
	aiReg := NewAIContractRegistry(reg)
	addr, err := aiReg.DeployAIContract([]byte{0x01}, "model", "", 5, "owner")
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	if h, ok := aiReg.ModelHash(addr); !ok || h != "model" {
		t.Fatalf("model hash mismatch")
	}
	if _, _, err := aiReg.InvokeAIContract(addr, []byte("in"), 5); err != nil {
		t.Fatalf("invoke: %v", err)
	}
}
