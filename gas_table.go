package synnergy

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	ilog "synnergy/internal/log"
	"synnergy/internal/telemetry"
)

// GasTable maps opcode names to their base gas cost. Stage 12 introduces pricing
// for data distribution monitor queries and Stage 13 adds initial costs for DEX
// liquidity lookups so dashboards can surface fees. Stage 23 registers consensus
// and governance CLI operations so their fees are visible to users. Stage 24
// extends coverage to cross-chain bridges, protocol registration and Plasma
// management so inter-network workflows remain predictable. Stage 42 records
// lock-mint and burn-release transfer opcodes so cross-chain transaction relays
// surface their costs alongside bridge operations. Stage 25 adds node
// management operations (full, light, mining, staking, watchtower and warfare
// nodes) so dashboards can price infrastructure actions. Stage 22 prices node status queries for the node operations dashboard. Stage 29 introduces
// pricing for deployable smart contract templates including token faucets,
// storage markets, DAO governance, NFT minting and AI model exchanges. Stage 35
// records storage marketplace operations so listing and deal workflows are gas
// priced. Stage 36 adds NFT marketplace operations to price minting and trading
// tokens. Stage 38 prices biometric security node enrollment and authentication
// commands so secure transaction flows surface their costs. Stage 39 records
// liquidity view operations for the DEX screener. Stage 40 adds monetary coin
// queries and compliance management costs so wallets can surface the price of
// reward, supply and KYC operations. Stage 46 prices ledger, light node,
// liquidity pool and loan pool operations so their costs are exposed via the CLI.
// Stage 59 registers content node management and content storage operations so
// hosts and registries expose deterministic pricing for publish, retrieval and
// discovery workflows used by the CLI and web interfaces. Stage 65 records DAO
// role update operations so governance tooling prices admin-managed role
// changes and authority term renewals. Stage 73 adds regulatory node approval,
// flagging, log query and audit operations so compliance tooling surfaces predictable
// gas costs. Stage 80 introduces Synthron treasury orchestration so monetary
// policy controls expose deterministic costs for issuing, burning and reconciling
// supply across the VM, consensus and authority layers.
type GasTable map[string]uint64

// GasMetadata captures descriptive information about opcodes so enterprise
// tooling can surface user friendly pricing insights without scraping the
// documentation at runtime.
type GasMetadata struct {
	Name        string `json:"name"`
	Cost        uint64 `json:"cost"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// DefaultGasCost is used when an opcode is missing from the guide.
// A low value keeps experimental opcodes affordable during development.
const DefaultGasCost uint64 = 1

var (
	gasOnce       sync.Once
	gasCache      GasTable
	gasMu         sync.RWMutex
	metadataCache map[string]GasMetadata
	metadataMu    sync.RWMutex
)

// ErrInvalidGasRegistration is returned when registering an opcode with an empty
// name or zero cost.
var ErrInvalidGasRegistration = errors.New("invalid gas registration")

// loadGasTable parses gas_table_list.md and caches the result.
func loadGasTable() {
	_, span := telemetry.Tracer("gas_table").Start(context.Background(), "GasTable.load")
	defer span.End()

	tbl := make(GasTable)
	meta := make(map[string]GasMetadata)
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(filename), "docs", "reference", "gas_table_list.md")
	data, err := os.ReadFile(path)
	if err != nil {
		span.RecordError(err)
		ilog.Error("gas_table_load", "error", err)
		gasCache = tbl
		metadataMu.Lock()
		metadataCache = meta
		metadataMu.Unlock()
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
		cost, parseErr := strconv.ParseUint(costStr, 10, 64)
		if parseErr != nil {
			cost = DefaultGasCost
			ilog.Error("gas_table_parse", "opcode", name, "error", parseErr)
		}
		tbl[name] = cost
		meta[name] = GasMetadata{
			Name:        name,
			Cost:        cost,
			Category:    "documentation",
			Description: "Documented reference schedule entry",
		}
	}
	if err := scanner.Err(); err != nil {
		span.RecordError(err)
		ilog.Error("gas_table_scan", "error", err)
	}
	gasCache = tbl
	metadataMu.Lock()
	metadataCache = meta
	metadataMu.Unlock()
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

// MustGasCost returns the gas price for an opcode and panics if it is missing.
// It is useful during initialization of critical modules where undefined
// pricing would indicate a misconfigured build or documentation drift.
func MustGasCost(opcode string) uint64 {
	if c, ok := LoadGasTable()[opcode]; ok {
		return c
	}
	panic(fmt.Sprintf("missing gas cost for opcode %s", opcode))
}

// HasOpcode reports whether a gas price is defined for the opcode.
func HasOpcode(name string) bool {
	tbl := LoadGasTable()
	_, ok := tbl[name]
	return ok
}

// EnsureGasSchedule guarantees that each opcode in schedule has a registered gas
// price. New entries are registered with the provided cost while existing
// entries are updated if the documented price has changed. The function returns
// the list of opcodes that were inserted for the first time, enabling callers to
// surface enterprise readiness checks across the CLI and web interfaces.
func EnsureGasSchedule(schedule map[string]uint64) ([]string, error) {
	if len(schedule) == 0 {
		return nil, nil
	}

	LoadGasTable()

	inserted := make([]string, 0, len(schedule))
	for name, cost := range schedule {
		if name == "" {
			return inserted, fmt.Errorf("%w: name", ErrInvalidGasRegistration)
		}
		if cost == 0 {
			return inserted, fmt.Errorf("%w: cost", ErrInvalidGasRegistration)
		}
		if !HasOpcode(name) {
			inserted = append(inserted, name)
		}
		if err := RegisterGasCost(name, cost); err != nil {
			return inserted, err
		}
		metadataMu.Lock()
		entry := metadataCache[name]
		entry.Name = name
		entry.Cost = cost
		if entry.Category == "" {
			entry.Category = "schedule"
		}
		if entry.Description == "" {
			entry.Description = "Registered via EnsureGasSchedule"
		}
		metadataCache[name] = entry
		metadataMu.Unlock()
	}

	if len(inserted) > 1 {
		sort.Strings(inserted)
	}
	return inserted, nil
}

// RegisterGasCost allows the CLI or tests to inject additional opcode pricing
// at runtime. It is safe for concurrent use and validates input to prevent
// accidental misconfiguration.
func RegisterGasCost(name string, cost uint64) error {
	if name == "" {
		return fmt.Errorf("%w: name", ErrInvalidGasRegistration)
	}
	if cost == 0 {
		return fmt.Errorf("%w: cost", ErrInvalidGasRegistration)
	}
	gasMu.Lock()
	defer gasMu.Unlock()
	if gasCache == nil {
		gasCache = make(GasTable)
	}
	gasCache[name] = cost
	metadataMu.Lock()
	if metadataCache == nil {
		metadataCache = make(map[string]GasMetadata)
	}
	entry := metadataCache[name]
	entry.Name = name
	entry.Cost = cost
	if entry.Category == "" {
		entry.Category = "runtime"
	}
	if entry.Description == "" {
		entry.Description = "Registered dynamically"
	}
	metadataCache[name] = entry
	metadataMu.Unlock()
	return nil
}

// RegisterGasMetadata enriches an opcode registration with category and
// descriptive text. It is used by bootstrap workflows so enterprise clients can
// surface contextual pricing data in dashboards and tests.
func RegisterGasMetadata(name string, cost uint64, category, description string) error {
	if err := RegisterGasCost(name, cost); err != nil {
		return err
	}
	if category == "" {
		category = "runtime"
	}
	metadataMu.Lock()
	entry := metadataCache[name]
	entry.Name = name
	entry.Cost = cost
	entry.Category = category
	entry.Description = description
	metadataCache[name] = entry
	metadataMu.Unlock()
	return nil
}

// GasMetadataFor retrieves the metadata for a specific opcode. The boolean
// indicates whether metadata exists for the key.
func GasMetadataFor(name string) (GasMetadata, bool) {
	LoadGasTable()
	metadataMu.RLock()
	defer metadataMu.RUnlock()
	entry, ok := metadataCache[name]
	if !ok {
		return GasMetadata{}, false
	}
	return entry, true
}

// GasCatalogue returns a deterministic list of metadata entries suitable for
// CLI rendering and telemetry exports.
func GasCatalogue() []GasMetadata {
	LoadGasTable()
	metadataMu.RLock()
	entries := make([]GasMetadata, 0, len(metadataCache))
	for _, entry := range metadataCache {
		entries = append(entries, entry)
	}
	metadataMu.RUnlock()
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	return entries
}

// ResetGasTable clears the cached table. Primarily used in tests to reload
// updated pricing without restarting the process.
func ResetGasTable() {
	gasMu.Lock()
	gasCache = nil
	gasMu.Unlock()
	metadataMu.Lock()
	metadataCache = nil
	metadataMu.Unlock()
	gasOnce = sync.Once{}
}

// MissingOpcodes reports which names are absent from the opcode catalogue. The
// helper is primarily used by diagnostics exposed through the CLI and function
// web to ensure new runtime features are documented.
func MissingOpcodes(names []string) []string {
	missing := make([]string, 0, len(names))
	for _, name := range names {
		if SNVMOpcodeByName(name) == 0 {
			missing = append(missing, name)
		}
	}
	if len(missing) > 1 {
		sort.Strings(missing)
	}
	return missing
}
