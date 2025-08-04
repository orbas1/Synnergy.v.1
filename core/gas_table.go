package core

// SetGasCost updates the gas cost for a specific opcode.
// If the global gas table has not been initialised yet it
// will be reset to the default values before applying the
// override.
func SetGasCost(op Opcode, cost uint64) {
	if gasTable == nil {
		initGasTable()
	}
	gasTable[op] = cost
}

// ResetGasTable restores the global gas table to the default
// values returned by DefaultGasTable. This provides a simple
// way for governance logic or tests to ensure a known set of
// gas prices.
func ResetGasTable() {
	initGasTable()
}
