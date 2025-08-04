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
