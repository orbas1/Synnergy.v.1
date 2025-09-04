package synnergy

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	ilog "synnergy/internal/log"
	"synnergy/internal/telemetry"
)

// GasTable maps opcode names to their base gas cost. Stage 12 introduces pricing for data distribution monitor queries. Stage 23 registers
// consensus and governance CLI operations so their fees are visible to users.
// Stage 24 extends coverage to cross-chain bridges, protocol registration and
// Plasma management so inter-network workflows remain predictable. Stage 25
// adds node management operations (full, light, mining, staking, watchtower and
// warfare nodes) so dashboards can price infrastructure actions. Stage 29 introduces pricing for deployable smart contract templates including token faucets, storage markets, DAO governance, NFT minting and AI model exchanges. Stage 35 records storage marketplace operations so listing and deal workflows are gas priced. Stage 36 adds NFT marketplace operations to price minting and trading tokens. Stage 39 records liquidity view operations for the DEX screener.
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
	_, span := telemetry.Tracer().Start(context.Background(), "GasTable.load")
	defer span.End()

	tbl := make(GasTable)
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(filename), "docs", "reference", "gas_table_list.md")
	data, err := os.ReadFile(path)
	if err != nil {
		span.RecordError(err)
		ilog.Error("gas_table_load", "error", err)
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
			ilog.Error("gas_table_parse", "opcode", name)
		}
	}
	if err := scanner.Err(); err != nil {
		span.RecordError(err)
		ilog.Error("gas_table_scan", "error", err)
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

// HasOpcode reports whether a gas price is defined for the opcode.
func HasOpcode(name string) bool {
	tbl := LoadGasTable()
	_, ok := tbl[name]
	return ok
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

// ResetGasTable clears the cached table. Primarily used in tests to reload
// updated pricing without restarting the process.
func ResetGasTable() {
	gasMu.Lock()
	gasCache = nil
	gasMu.Unlock()
	gasOnce = sync.Once{}
}
