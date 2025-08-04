package core

// GasTable defines the gas cost for each opcode.
type GasTable map[Opcode]uint64

// gasTable holds the active gas pricing used by the VM. It is initialised
// from DefaultGasTable and can be modified at runtime via SetGasCost.
var gasTable = DefaultGasTable()

// GasCost returns the gas cost for the given opcode. If an opcode is missing
// from the table it returns zero, allowing callers to treat unpriced opcodes
// as free but still present in the catalogue.
func GasCost(op Opcode) uint64 {
	return gasTable[op]
}

// initGasTable resets the global gas table to defaults. It is invoked during
// opcode initialisation once the catalogue is populated.
func initGasTable() {
	gasTable = DefaultGasTable()
}
