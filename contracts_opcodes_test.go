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
		{"OpLockAndMint", OpLockAndMint, 0x090004},
		{"OpBurnAndRelease", OpBurnAndRelease, 0x090005},
		{"OpLiquidityAdd", OpLiquidityAdd, 0x0F0004},
		{"OpQueryOracle", OpQueryOracle, 0x0A0008},
		{"OpMintToken", OpMintToken, 0x0E0014},
		{"OpMultiSigSubmit", OpMultiSigSubmit, 0x200001},
		{"OpMultiSigConfirm", OpMultiSigConfirm, 0x200002},
		{"OpMultiSigRevoke", OpMultiSigRevoke, 0x200003},
		{"OpMultiSigExecute", OpMultiSigExecute, 0x200004},
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
		OpLockAndMint,
		OpBurnAndRelease,
		OpLiquidityAdd,
		OpQueryOracle,
		OpMintToken,
		OpMultiSigSubmit,
		OpMultiSigConfirm,
		OpMultiSigRevoke,
		OpMultiSigExecute,
	}

	seen := make(map[uint32]struct{}, len(ops))
	for _, op := range ops {
		if _, ok := seen[op]; ok {
			t.Fatalf("duplicate opcode value %#x", op)
		}
		seen[op] = struct{}{}
	}
}
