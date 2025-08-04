package core

// GasTable defines the gas cost for each opcode.
type GasTable map[Opcode]uint64

// DefaultGasTable returns a basic gas table.
func DefaultGasTable() GasTable {
	return GasTable{
		OpNoop:     1,
		OpTransfer: 10,
	}
}

var gasTable = DefaultGasTable()

// GasCost returns the gas cost for the given opcode.
func GasCost(op Opcode) uint64 {
	return gasTable[op]
}

// initGasTable resets the global gas table to defaults.
func initGasTable() {
	gasTable = DefaultGasTable()
}
