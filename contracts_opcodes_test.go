package synnergy

import "testing"

// TestContractOpcodesValues ensures that all contract opcodes remain stable
// and match their documented hexadecimal representations.
func TestContractOpcodesValues(t *testing.T) {
	tests := []struct {
		name string
		op   uint32
		want uint32
	}{
		{"OpInitContracts", OpInitContracts, 0x080001},
		{"OpPauseContract", OpPauseContract, 0x080006},
		{"OpResumeContract", OpResumeContract, 0x080007},
		{"OpUpgradeContract", OpUpgradeContract, 0x080008},
		{"OpContractInfo", OpContractInfo, 0x080009},
		{"OpDeployAIContract", OpDeployAIContract, 0x010001},
		{"OpInvokeAIContract", OpInvokeAIContract, 0x010002},
	}
	for _, tt := range tests {
		if tt.op != tt.want {
			t.Errorf("%s = %#x, want %#x", tt.name, tt.op, tt.want)
		}
		if tt.op == 0 {
			t.Errorf("%s opcode must be non-zero", tt.name)
		}
	}
}

// TestContractOpcodesUnique verifies the opcodes are unique.
func TestContractOpcodesUnique(t *testing.T) {
	ops := []uint32{
		OpInitContracts,
		OpPauseContract,
		OpResumeContract,
		OpUpgradeContract,
		OpContractInfo,
		OpDeployAIContract,
		OpInvokeAIContract,
	}

	seen := make(map[uint32]struct{}, len(ops))
	for _, op := range ops {
		if _, ok := seen[op]; ok {
			t.Fatalf("duplicate opcode value %#x", op)
		}
		seen[op] = struct{}{}
	}
}
