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
	"strconv"
	"strings"
	"sync"
	"time"

	ilog "synnergy/internal/log"
	"synnergy/internal/telemetry"
)

// GasTable maps opcode names to their base gas cost. Stage 12 introduces pricing
// for data distribution monitor queries and Stage 13 adds initial costs for DEX
// liquidity lookups so dashboards can surface fees. Stage 23 registers consensus
// and governance CLI operations so their fees are visible to users. Stage 24
// extends coverage to cross-chain bridges, protocol registration and Plasma
// management so inter-network workflows remain predictable. Stage 42 records
// lock-mint and burn-release transfer operations so cross-chain transaction relays
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
// gas costs.
type GasTable map[string]uint64

// GasSnapshot represents the immutable state of the gas registry at a point in
// time. Consumers subscribe to snapshots to keep long-lived services (CLI,
// authority nodes, web tier) synchronised without re-reading documentation on
// every lookup.
type GasSnapshot struct {
	Table    GasTable
	Version  uint64
	LoadedAt time.Time
}

// DefaultGasCost is used when an opcode is missing from the guide.
// A low value keeps experimental opcodes affordable during development.
const DefaultGasCost uint64 = 1

var (
	gasOnce  sync.Once
	gasCache GasTable
	gasMu    sync.RWMutex
	gasVer   uint64
	gasTime  time.Time

	gasSubscribersMu sync.RWMutex
	gasSubscribers   map[chan GasSnapshot]struct{}
)

// ErrInvalidGasRegistration is returned when registering an opcode with an empty
// name or zero cost.
var ErrInvalidGasRegistration = errors.New("invalid gas registration")

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
		gasMu.Lock()
		gasCache = tbl
		gasVer++
		gasTime = time.Now().UTC()
		snapshot := snapshotLocked()
		gasMu.Unlock()
		notifyGasSubscribers(snapshot)
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
	gasMu.Lock()
	if gasCache == nil {
		gasCache = tbl
	} else {
		for name, cost := range tbl {
			if _, exists := gasCache[name]; !exists {
				gasCache[name] = cost
			}
		}
	}
	gasVer++
	gasTime = time.Now().UTC()
	snapshot := snapshotLocked()
	gasMu.Unlock()
	notifyGasSubscribers(snapshot)
}

// LoadGasTable returns the cached gas table, loading it on first use.
func LoadGasTable() GasTable {
	gasOnce.Do(loadGasTable)
	gasMu.RLock()
	defer gasMu.RUnlock()
	return cloneGasTable(gasCache)
}

// GasCost returns the gas price for a given opcode name. If the opcode is not
// present in the table, DefaultGasCost is returned.
func GasCost(opcode string) uint64 {
	gasOnce.Do(loadGasTable)
	gasMu.RLock()
	defer gasMu.RUnlock()
	if c, ok := gasCache[opcode]; ok {
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
	gasOnce.Do(loadGasTable)
	gasMu.RLock()
	defer gasMu.RUnlock()
	_, ok := gasCache[name]
	return ok
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
	if gasCache == nil {
		gasCache = make(GasTable)
	}
	unchanged := gasCache[name] == cost
	gasCache[name] = cost
	gasVer++
	gasTime = time.Now().UTC()
	snapshot := snapshotLocked()
	gasMu.Unlock()
	if !unchanged {
		notifyGasSubscribers(snapshot)
	}
	return nil
}

// ResetGasTable clears the cached table. Primarily used in tests to reload
// updated pricing without restarting the process.
func ResetGasTable() {
	gasMu.Lock()
	gasCache = nil
	gasVer++
	gasTime = time.Now().UTC()
	snapshot := snapshotLocked()
	gasMu.Unlock()
	gasOnce = sync.Once{}
	notifyGasSubscribers(snapshot)
}

// ReloadGasTable invalidates the cache and loads a fresh copy from disk. The
// returned snapshot mirrors the state observed by subscribers after the
// reload.
func ReloadGasTable() GasSnapshot {
	gasMu.Lock()
	overrides := cloneGasTable(gasCache)
	gasMu.Unlock()

	gasOnce = sync.Once{}
	gasMu.Lock()
	gasCache = overrides
	gasMu.Unlock()

	gasOnce.Do(loadGasTable)
	return SnapshotGasTable()
}

// SnapshotGasTable returns the current gas configuration along with versioning
// metadata so callers can detect staleness.
func SnapshotGasTable() GasSnapshot {
	gasOnce.Do(loadGasTable)
	gasMu.RLock()
	defer gasMu.RUnlock()
	return snapshotLocked()
}

// SubscribeGasTable streams immutable snapshots of the gas table. The caller
// receives the current snapshot immediately and subsequent updates whenever the
// registry changes. The returned cancel function must be invoked to stop the
// stream and prevent goroutine leaks.
func SubscribeGasTable(buffer int) (<-chan GasSnapshot, func()) {
	if buffer <= 0 {
		buffer = 1
	}
	gasOnce.Do(loadGasTable)
	ch := make(chan GasSnapshot, buffer)
	gasSubscribersMu.Lock()
	if gasSubscribers == nil {
		gasSubscribers = make(map[chan GasSnapshot]struct{})
	}
	gasSubscribers[ch] = struct{}{}
	gasSubscribersMu.Unlock()
	ch <- SnapshotGasTable()
	var once sync.Once
	cancel := func() {
		once.Do(func() {
			gasSubscribersMu.Lock()
			if _, ok := gasSubscribers[ch]; ok {
				delete(gasSubscribers, ch)
				close(ch)
			}
			gasSubscribersMu.Unlock()
		})
	}
	return ch, cancel
}

func snapshotLocked() GasSnapshot {
	return GasSnapshot{Table: cloneGasTable(gasCache), Version: gasVer, LoadedAt: gasTime}
}

func cloneGasTable(src GasTable) GasTable {
	clone := make(GasTable, len(src))
	for k, v := range src {
		clone[k] = v
	}
	return clone
}

func notifyGasSubscribers(snapshot GasSnapshot) {
	gasSubscribersMu.RLock()
	defer gasSubscribersMu.RUnlock()
	for ch := range gasSubscribers {
		select {
		case ch <- snapshot:
		default:
		}
	}
}
