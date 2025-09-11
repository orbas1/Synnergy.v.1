package core

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
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
	applyEnvOverrides(tbl)
	return tbl
}

// parseGasGuide reads gas_table_list.md and extracts price overrides.
// The file is expected to contain markdown tables with backtick-quoted opcode
// names and numeric gas costs.
func parseGasGuide() map[string]uint64 {
	_, filename, _, _ := runtime.Caller(0)
	coreDir := filepath.Dir(filename)
	rootDir := filepath.Dir(coreDir)
	path := filepath.Clean(filepath.Join(rootDir, "docs", "reference", "gas_table_list.md"))
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

// applyEnvOverrides allows operators to supply a comma-separated list of
// `Name=Cost` pairs via the SYN_GAS_OVERRIDES environment variable. Missing or
// malformed entries are ignored. Names are resolved through the opcode
// catalogue so overrides only apply to recognised operations.
func applyEnvOverrides(tbl GasTable) {
	overrides := os.Getenv("SYN_GAS_OVERRIDES")
	if overrides == "" {
		return
	}
	for _, p := range strings.Split(overrides, ",") {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) != 2 {
			continue
		}
		cost, err := strconv.ParseUint(kv[1], 10, 64)
		if err != nil {
			continue
		}
		if op, ok := nameToOp[kv[0]]; ok {
			tbl[op] = cost
		}
	}
}

// SetGasCost updates the gas cost for a specific opcode at runtime. This allows
// governance or tests to tweak opcode pricing without rebuilding the binary.
func SetGasCost(op Opcode, cost uint64) {
	gasMu.Lock()
	gasTable[op] = cost
	gasMu.Unlock()
}

// GasTableSnapshot returns a copy of the current gas table. The snapshot can be
// used for inspection or persistence.
func GasTableSnapshot() GasTable {
	gasMu.RLock()
	defer gasMu.RUnlock()
	snapshot := make(GasTable, len(gasTable))
	for op, c := range gasTable {
		snapshot[op] = c
	}
	return snapshot
}

// GasTableSnapshotJSON serialises the current gas schedule to JSON. Opcode keys
// are rendered in hexadecimal ("0x000000") form so external tooling does not
// need awareness of internal catalogue names. Map iteration in Go is randomised,
// so the function sorts the keys to produce deterministic output for hashing or
// signature checks.
func GasTableSnapshotJSON() ([]byte, error) {
	snap := GasTableSnapshot()
	type kv struct {
		key  string
		cost uint64
	}
	pairs := make([]kv, 0, len(snap))
	for op, cost := range snap {
		pairs = append(pairs, kv{fmt.Sprintf("0x%06X", op), cost})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].key < pairs[j].key })
	var b bytes.Buffer
	b.WriteByte('{')
	for i, p := range pairs {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%q:%d", p.key, p.cost)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

// WriteGasTableSnapshot persists the current gas table snapshot to the
// provided path. The file is written in the same deterministic JSON format as
// returned by GasTableSnapshotJSON.
func WriteGasTableSnapshot(path string) error {
	data, err := GasTableSnapshotJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// GasCostByName returns the gas price for an exported function name. It
// resolves the name through the opcode catalogue and falls back to
// DefaultGasCost when the function or its gas entry is unknown.
func GasCostByName(name string) uint64 {
	mu.RLock()
	op, ok := nameToOp[name]
	mu.RUnlock()
	if !ok {
		return DefaultGasCost
	}
	gasMu.RLock()
	cost, ok := gasTable[op]
	gasMu.RUnlock()
	if !ok {
		return DefaultGasCost
	}
	return cost
}
