package core

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	synn "synnergy"
	"synnergy/internal/telemetry"
)

// EnterpriseOption configures the enterprise orchestrator at construction time.
type EnterpriseOption func(*EnterpriseOrchestrator)

// WithGasSchedule allows callers to augment the default opcode gas schedule
// enforced by the orchestrator. Costs provided here are merged with the
// built-in Stage 78 defaults so operators can tune deployments without editing
// source code.
func WithGasSchedule(schedule map[string]uint64) EnterpriseOption {
	return func(o *EnterpriseOrchestrator) {
		if len(schedule) == 0 {
			return
		}
		if o.gas == nil {
			o.gas = make(map[string]uint64, len(schedule))
		}
		for name, cost := range schedule {
			o.gas[name] = cost
		}
	}
}

// EnterpriseDiagnostics summarises health across the virtual machine, consensus
// mesh, wallet and authority node registry. Results are cached for short periods
// so repeated CLI or web requests remain inexpensive.
type EnterpriseDiagnostics struct {
	Timestamp         time.Time         `json:"timestamp"`
	VMRunning         bool              `json:"vmRunning"`
	VMMode            string            `json:"vmMode"`
	VMConcurrency     int               `json:"vmConcurrency"`
	ConsensusNetworks int               `json:"consensusNetworks"`
	ConsensusRelayers int               `json:"consensusRelayers"`
	AuthorityNodes    int               `json:"authorityNodes"`
	AuthorityRoles    map[string]int    `json:"authorityRoles"`
	WalletAddress     string            `json:"walletAddress"`
	WalletSealed      bool              `json:"walletSealed"`
	LedgerHeight      int               `json:"ledgerHeight"`
	GasCoverage       map[string]uint64 `json:"gasCoverage"`
	MissingOpcodes    []string          `json:"missingOpcodes"`
	InsertedOpcodes   []string          `json:"insertedOpcodes,omitempty"`
	GasLastSyncedAt   time.Time         `json:"gasLastSyncedAt"`
}

// EnterpriseOrchestrator coordinates Stage 78 subsystems so enterprise operators
// can verify readiness from the CLI, automated tests or the function web.
type EnterpriseOrchestrator struct {
	mu                sync.RWMutex
	vm                *SimpleVM
	consensus         *ConsensusNetworkManager
	wallet            *Wallet
	registry          *AuthorityNodeRegistry
	ledger            *Ledger
	gas               map[string]uint64
	last              EnterpriseDiagnostics
	pendingInsertions []string
	lastSync          time.Time
}

// NewEnterpriseOrchestrator boots a heavy VM, authorises a consensus relayer,
// registers the orchestrator wallet as an authority node and validates Stage 78
// opcode pricing. The resulting instance exposes diagnostics and helper
// functions to wire additional subsystems.
func NewEnterpriseOrchestrator(ctx context.Context, opts ...EnterpriseOption) (*EnterpriseOrchestrator, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	vm := NewSimpleVM(VMHeavy)
	if err := vm.Start(); err != nil {
		return nil, fmt.Errorf("start vm: %w", err)
	}

	wallet, err := NewWallet()
	if err != nil {
		return nil, fmt.Errorf("wallet init: %w", err)
	}

	orchestrator := &EnterpriseOrchestrator{
		vm:        vm,
		consensus: NewConsensusNetworkManager(),
		wallet:    wallet,
		registry:  NewAuthorityNodeRegistry(),
		ledger:    NewLedger(),
		gas: map[string]uint64{
			"EnterpriseBootstrap":      120,
			"EnterpriseConsensusSync":  95,
			"EnterpriseWalletSeal":     60,
			"EnterpriseNodeAudit":      75,
			"EnterpriseAuthorityElect": 80,
		},
	}

	for _, opt := range opts {
		opt(orchestrator)
	}

	inserted, err := synn.EnsureGasSchedule(orchestrator.gas)
	if err != nil {
		return nil, fmt.Errorf("ensure gas schedule: %w", err)
	}
	orchestrator.pendingInsertions = inserted
	orchestrator.lastSync = time.Now().UTC()

	orchestrator.consensus.AuthorizeRelayer(orchestrator.wallet.Address)
	if !orchestrator.registry.IsAuthorityNode(orchestrator.wallet.Address) {
		if _, err := orchestrator.registry.Register(orchestrator.wallet.Address, "orchestrator"); err != nil {
			return nil, fmt.Errorf("register orchestrator authority: %w", err)
		}
	}

	orchestrator.refreshDiagnostics(ctx)
	return orchestrator, nil
}

// Diagnostics returns cached health information or refreshes it if the previous
// snapshot is older than one second. The behaviour matches CLI and web polling
// intervals so repeated checks do not thrash shared state.
func (o *EnterpriseOrchestrator) Diagnostics(ctx context.Context) EnterpriseDiagnostics {
	o.mu.RLock()
	diag := o.last
	o.mu.RUnlock()
	if time.Since(diag.Timestamp) < time.Second {
		return diag
	}
	return o.refreshDiagnostics(ctx)
}

// RegisterConsensusNetwork provisions a new cross-consensus connection using the
// orchestrator wallet as the authorised relayer. Diagnostics are refreshed so
// downstream dashboards immediately reflect the topology change.
func (o *EnterpriseOrchestrator) RegisterConsensusNetwork(ctx context.Context, source, target string) (int, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	_, span := telemetry.Tracer("core.enterprise").Start(ctx, "EnterpriseOrchestrator.RegisterConsensusNetwork")
	defer span.End()

	id, err := o.consensus.RegisterNetwork(source, target, o.wallet.Address)
	if err != nil {
		span.RecordError(err)
		return 0, err
	}
	o.refreshDiagnostics(ctx)
	return id, nil
}

// RegisterAuthorityNode adds an authority node to the shared registry. When
// addr is empty the orchestrator wallet acts as the operator identity.
func (o *EnterpriseOrchestrator) RegisterAuthorityNode(ctx context.Context, addr, role string) (*AuthorityNode, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if addr == "" {
		addr = o.wallet.Address
	}
	node, err := o.registry.Register(addr, role)
	if err != nil {
		return nil, err
	}
	o.refreshDiagnostics(ctx)
	return node, nil
}

// EnterpriseBootstrap ensures every Stage 78 subsystem is initialised and returns
// an updated diagnostics snapshot. It is safe to invoke multiple times.
func (o *EnterpriseOrchestrator) EnterpriseBootstrap(ctx context.Context) (EnterpriseDiagnostics, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := o.vm.Start(); err != nil {
		return EnterpriseDiagnostics{}, fmt.Errorf("vm start: %w", err)
	}
	if _, err := o.EnterpriseWalletSeal(); err != nil {
		return EnterpriseDiagnostics{}, err
	}
	if err := o.EnterpriseConsensusSync(ctx); err != nil {
		return EnterpriseDiagnostics{}, err
	}
	if _, err := o.EnterpriseAuthorityElect(ctx, "orchestrator"); err != nil {
		return EnterpriseDiagnostics{}, err
	}
	if err := o.EnterpriseNodeAudit(ctx); err != nil {
		return EnterpriseDiagnostics{}, err
	}
	return o.refreshDiagnostics(ctx), nil
}

// EnterpriseConsensusSync ensures at least one consensus bridging network is
// configured using the orchestrator wallet as the authorised relayer.
func (o *EnterpriseOrchestrator) EnterpriseConsensusSync(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(o.consensus.ListNetworks()) > 0 {
		return nil
	}
	addr, err := o.EnterpriseWalletSeal()
	if err != nil {
		return err
	}
	if _, err := o.consensus.RegisterNetwork("ProofOfStake", "Synnergy-PBFT", addr); err != nil {
		return fmt.Errorf("register consensus network: %w", err)
	}
	o.refreshDiagnostics(ctx)
	return nil
}

// EnterpriseWalletSeal validates that the orchestrator wallet is available and
// returns its address.
func (o *EnterpriseOrchestrator) EnterpriseWalletSeal() (string, error) {
	if o.wallet == nil {
		return "", errors.New("wallet unavailable")
	}
	if o.wallet.Address == "" {
		return "", errors.New("wallet address not initialised")
	}
	return o.wallet.Address, nil
}

// EnterpriseAuthorityElect ensures the orchestrator wallet is registered as an
// authority node under the provided role.
func (o *EnterpriseOrchestrator) EnterpriseAuthorityElect(ctx context.Context, role string) (*AuthorityNode, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	addr, err := o.EnterpriseWalletSeal()
	if err != nil {
		return nil, err
	}
	if role == "" {
		role = "orchestrator"
	}
	if o.registry.IsAuthorityNode(addr) {
		node, infoErr := o.registry.Info(addr)
		if infoErr != nil {
			return nil, infoErr
		}
		o.refreshDiagnostics(ctx)
		return node, nil
	}
	node, err := o.registry.Register(addr, role)
	if err != nil {
		return nil, err
	}
	o.refreshDiagnostics(ctx)
	return node, nil
}

// EnterpriseNodeAudit touches the ledger to confirm read access is healthy and
// refreshes diagnostics for monitoring pipelines.
func (o *EnterpriseOrchestrator) EnterpriseNodeAudit(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if o.ledger == nil {
		return errors.New("ledger unavailable")
	}
	_, _ = o.ledger.Head()
	o.refreshDiagnostics(ctx)
	return nil
}

// SyncGasSchedule merges schedule with the orchestrator baseline and ensures the
// global gas table is updated. Diagnostics are returned so callers can verify
// pricing as part of automation pipelines.
func (o *EnterpriseOrchestrator) SyncGasSchedule(ctx context.Context, schedule map[string]uint64) (EnterpriseDiagnostics, error) {
	if len(schedule) > 0 {
		o.mu.Lock()
		if o.gas == nil {
			o.gas = make(map[string]uint64, len(schedule))
		}
		for name, cost := range schedule {
			o.gas[name] = cost
		}
		o.mu.Unlock()
		inserted, err := synn.EnsureGasSchedule(o.gas)
		if err != nil {
			return EnterpriseDiagnostics{}, err
		}
		o.mu.Lock()
		o.pendingInsertions = inserted
		o.lastSync = time.Now().UTC()
		o.mu.Unlock()
	}
	return o.refreshDiagnostics(ctx), nil
}

// refreshDiagnostics recomputes the diagnostics snapshot and records telemetry
// for observability tooling used by the function web.
func (o *EnterpriseOrchestrator) refreshDiagnostics(ctx context.Context) EnterpriseDiagnostics {
	if ctx == nil {
		ctx = context.Background()
	}
	_, span := telemetry.Tracer("core.enterprise").Start(ctx, "EnterpriseOrchestrator.refreshDiagnostics")
	defer span.End()

	networks := o.consensus.ListNetworks()
	relayers := o.consensus.AuthorizedRelayers()
	nodes := o.registry.List()

	walletAddress := ""
	walletSealed := false
	if o.wallet != nil {
		walletAddress = o.wallet.Address
		walletSealed = walletAddress != ""
	}

	diag := EnterpriseDiagnostics{
		Timestamp:         time.Now().UTC(),
		VMRunning:         o.vm.Status(),
		VMMode:            o.vm.Mode().String(),
		VMConcurrency:     o.vm.Concurrency(),
		ConsensusNetworks: len(networks),
		ConsensusRelayers: len(relayers),
		AuthorityNodes:    len(nodes),
		AuthorityRoles:    make(map[string]int, len(nodes)),
		WalletAddress:     walletAddress,
		WalletSealed:      walletSealed,
		GasCoverage:       make(map[string]uint64, len(o.gas)),
		GasLastSyncedAt:   o.lastSync,
	}
	height, _ := o.ledger.Head()
	diag.LedgerHeight = height

	missing := make([]string, 0)
	for name, expected := range o.gas {
		cost := synn.GasCost(name)
		diag.GasCoverage[name] = cost
		if cost == synn.DefaultGasCost && expected != synn.DefaultGasCost && !synn.HasOpcode(name) {
			missing = append(missing, name)
		}
	}
	sort.Strings(missing)
	diag.MissingOpcodes = missing

	for _, node := range nodes {
		role := node.Role
		if role == "" {
			role = "unspecified"
		}
		diag.AuthorityRoles[role]++
	}

	o.mu.Lock()
	if len(o.pendingInsertions) > 0 {
		diag.InsertedOpcodes = append(diag.InsertedOpcodes, o.pendingInsertions...)
		o.pendingInsertions = nil
	}
	o.last = diag
	o.mu.Unlock()

	return diag
}
