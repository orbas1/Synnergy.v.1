package core

import "sync"

// GasTable defines the gas cost for each opcode.
type GasTable map[Opcode]uint64

var (
	gasTable = DefaultGasTable()
	gasMu    sync.RWMutex
)

// GasCost returns the gas cost for the given opcode. If an opcode is missing
// from the table it falls back to DefaultGasCost so callers receive a
// deterministic price even when the guide lacks an entry.
func GasCost(op Opcode) uint64 {
	gasMu.RLock()
	cost, ok := gasTable[op]
	gasMu.RUnlock()
	if !ok {
		return DefaultGasCost
	}
	return cost
}

// initGasTable resets the global gas table to defaults. It is invoked during
// opcode initialisation once the catalogue is populated.
func initGasTable() {
	gasMu.Lock()
	gasTable = DefaultGasTable()
	gasMu.Unlock()
}
