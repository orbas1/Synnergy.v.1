package core

import "testing"

func TestContractsOpcodes(t *testing.T) {
	names := []string{"InitContracts", "PauseContract", "ResumeContract", "UpgradeContract", "ContractInfo", "DeployAIContract", "InvokeAIContract"}
	seen := map[Opcode]bool{}
	for _, name := range names {
		op := opcodeByName(name)
		if op == 0 {
			t.Fatalf("opcode for %s missing", name)
		}
		if seen[op] {
			t.Fatalf("duplicate opcode %v", op)
		}
		seen[op] = true
	}
}
