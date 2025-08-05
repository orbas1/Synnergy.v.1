package core

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// DefaultGasCost is used for any opcode not explicitly priced in the guide.
const DefaultGasCost = 1

// DefaultGasTable builds a gas pricing table for all registered opcodes. It
// parses `gas_table_list.md` at runtime to pull concrete costs. Any
// opcode missing from the guide receives DefaultGasCost ensuring the table is
// exhaustive.
func DefaultGasTable() GasTable {
	overrides := parseGasGuide()
	tbl := make(GasTable, len(catalogue))
	for _, entry := range catalogue {
		if cost, ok := overrides[entry.name]; ok {
			tbl[entry.op] = cost
		} else {
			tbl[entry.op] = DefaultGasCost
		}
	}
	return tbl
}

// parseGasGuide reads gas_table_list.md and extracts price overrides.
// The file is expected to contain markdown tables with backtick-quoted opcode
// names and numeric gas costs.
func parseGasGuide() map[string]uint64 {
	_, filename, _, _ := runtime.Caller(0)
	coreDir := filepath.Dir(filename)
	rootDir := filepath.Dir(coreDir)
	path := filepath.Clean(filepath.Join(rootDir, "gas_table_list.md"))
	if !strings.HasPrefix(path, rootDir) {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	m := make(map[string]uint64)
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
			m[name] = cost
		}
	}
	return m
}

// SetGasCost updates the gas cost for a specific opcode at runtime. This allows
// governance or tests to tweak opcode pricing without rebuilding the binary.
func SetGasCost(op Opcode, cost uint64) {
	gasTable[op] = cost
}

// GasTableSnapshot returns a copy of the current gas table. The snapshot can be
// used for inspection or persistence.
func GasTableSnapshot() GasTable {
	snapshot := make(GasTable, len(gasTable))
	for op, c := range gasTable {
		snapshot[op] = c
	}
	return snapshot
}
