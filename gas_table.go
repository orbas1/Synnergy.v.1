package synnergy

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// GasTable maps opcode names to their base gas cost.
type GasTable map[string]uint64

// DefaultGasCost is used when an opcode is missing from the guide.
const DefaultGasCost uint64 = 100_000

// LoadGasTable parses opcode_and_gas_guide.md and returns a fully populated gas table.
// Lines in the guide are expected to contain markdown tables where the first column
// is the opcode name and the second column is the numeric gas cost. Any opcode not
// found in the guide receives DefaultGasCost.
func LoadGasTable() GasTable {
	tbl := make(GasTable)
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(filename), "opcode_and_gas_guide.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return tbl
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "| `") {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			continue
		}
		name := strings.Trim(parts[1], " `")
		costStr := strings.Trim(parts[2], " `")
		if cost, err := strconv.ParseUint(costStr, 10, 64); err == nil {
			tbl[name] = cost
		} else {
			tbl[name] = DefaultGasCost
		}
	}
	return tbl
}

// GasCost returns the gas price for a given opcode name. If the opcode is not
// present in the table, DefaultGasCost is returned.
func GasCost(opcode string) uint64 {
	tbl := LoadGasTable()
	if c, ok := tbl[opcode]; ok {
		return c
	}
	return DefaultGasCost
}
