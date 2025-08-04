package core

// contracts_opcodes.go exposes well-known opcode values for contract related
// operations. The opcodes are resolved at init time using the registry in
// opcode.go to ensure they stay consistent with the rest of the system.

// NOTE: these are variables rather than constants because the opcode catalogue
// is normalised during package initialisation. The helper ensures lookups panic
// if an opcode name is missing which would indicate a mismatch with the
// catalogue.

var (
	OpInitContracts    = opcodeByName("InitContracts")
	OpPauseContract    = opcodeByName("PauseContract")
	OpResumeContract   = opcodeByName("ResumeContract")
	OpUpgradeContract  = opcodeByName("UpgradeContract")
	OpContractInfo     = opcodeByName("ContractInfo")
	OpDeployAIContract = opcodeByName("DeployAIContract")
	OpInvokeAIContract = opcodeByName("InvokeAIContract")
)

func opcodeByName(name string) Opcode {
	b, err := ToBytecode(name)
	if err != nil {
		panic(err)
	}
	op, err := ParseOpcode(b)
	if err != nil {
		panic(err)
	}
	return op
}
