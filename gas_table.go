package synnergy

import (
	"bufio"
	"bytes"
	"log"
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
        gasMu    sync.RWMutex
)

// loadGasTable parses gas_table_list.md and caches the result.
func loadGasTable() {
	tbl := make(GasTable)
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(filename), "gas_table_list.md")
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("gas_table: %v", err)
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
        gasMu.RLock()
        defer gasMu.RUnlock()
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

// RegisterGasCost allows the CLI or tests to inject additional opcode pricing
// at runtime. It is safe for concurrent use.
func RegisterGasCost(name string, cost uint64) {
        gasMu.Lock()
        defer gasMu.Unlock()
        if gasCache == nil {
                gasCache = make(GasTable)
        }
        gasCache[name] = cost
}
