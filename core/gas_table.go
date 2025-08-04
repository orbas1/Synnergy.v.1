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

