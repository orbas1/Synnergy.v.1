package synnergy

// Contracts-related opcode identifiers. The numeric values are arbitrary but
// stable for testing and documentation purposes.
const (
	OpInitContracts    uint32 = 0x010000
	OpPauseContract    uint32 = 0x010001
	OpResumeContract   uint32 = 0x010002
	OpUpgradeContract  uint32 = 0x010003
	OpContractInfo     uint32 = 0x010004
	OpDeployAIContract uint32 = 0x010005
	OpInvokeAIContract uint32 = 0x010006
)
