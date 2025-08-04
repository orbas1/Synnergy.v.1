package core

// SetGasCost updates the gas cost for a specific opcode at runtime.
// This allows governance or tests to tweak opcode pricing without
// rebuilding the binary.
func SetGasCost(op Opcode, cost uint64) {
	gasTable[op] = cost
}

// GasTableSnapshot returns a copy of the current gas table. The
// snapshot can be used for inspection or persistence.
func GasTableSnapshot() GasTable {
	snapshot := make(GasTable, len(gasTable))
	for op, c := range gasTable {
		snapshot[op] = c
	}
	return snapshot
}
