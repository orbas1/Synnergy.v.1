package synnergy

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// GasTable maps opcode names to their base gas cost.
type GasTable map[string]uint64

// DefaultGasCost is used when an opcode is missing from the guide.
// A low value keeps experimental opcodes affordable during development.
const DefaultGasCost uint64 = 1

var (
	gasOnce  sync.Once
	gasCache GasTable
)

// loadGasTable parses gas_table_list.md and caches the result.
func loadGasTable() {
	tbl := make(GasTable)
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(filename), "gas_table_list.md")
	data, err := os.ReadFile(path)
	if err != nil {
		gasCache = tbl
		return
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
	gasCache = tbl
}

// LoadGasTable returns the cached gas table, loading it on first use.
func LoadGasTable() GasTable {
	gasOnce.Do(loadGasTable)
	return gasCache
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
