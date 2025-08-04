package core

// gasTable stores the global opcode gas cost table.
var gasTable GasTable

// initGasTable initialises the global gas table. It is invoked from opcode
// registration to ensure a baseline set of prices is available. The table can
// be updated at runtime via governance mechanisms.
func initGasTable() {
    if len(gasTable) == 0 {
        gasTable = DefaultGasTable()
    }
}

// Gas table placeholder.
// Future implementations should populate opcode gas costs.
// Placeholder file ensuring the core package builds while the full gas table
// implementation is developed elsewhere.

// gasTable holds the default gas costs for opcodes. It is initialised during
// package startup and may be extended in later stages.
var gasTable GasTable

// initGasTable populates the global gas table using DefaultGasTable. Called from
// the opcode dispatcher after all opcodes are registered.
func initGasTable() {
	gasTable = DefaultGasTable()
}

// GasCost returns the gas cost for a given opcode. Undefined opcodes return
// zero, allowing callers to treat them as invalid.
func GasCost(op Opcode) uint64 {
	if cost, ok := gasTable[op]; ok {
		return cost
	}
	return 0
}
