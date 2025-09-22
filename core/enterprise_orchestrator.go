package core

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	synn "synnergy"
	security "synnergy/internal/security"
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
	AuthorityNodes    int               `json:"authorityNodes"`
	WalletAddress     string            `json:"walletAddress"`
	LedgerHeight      int               `json:"ledgerHeight"`
	GasCoverage       map[string]uint64 `json:"gasCoverage"`
	MissingOpcodes    []string          `json:"missingOpcodes"`
	InsertedOpcodes   []string          `json:"insertedOpcodes,omitempty"`
	Resilience        *ResilienceReport `json:"resilience,omitempty"`
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
	failover          *FailoverManager
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

	consensusMgr := NewConsensusNetworkManager()
	registry := NewAuthorityNodeRegistry()
	ledger := NewLedger()
	failover := NewFailoverManager(wallet.Address, 10*time.Second,
		WithFailoverVirtualMachine(vm),
		WithFailoverConsensus(consensusMgr),
		WithFailoverWallet(wallet),
		WithFailoverRegistry(registry),
		WithFailoverLedger(ledger),
		WithFailoverSigner(security.NewKeyManager()),
	)
	failover.RegisterNode(FailoverNode{ID: wallet.Address, Role: "orchestrator", Region: "global", PublicKey: wallet.PublicKeyBytes()})

	orchestrator := &EnterpriseOrchestrator{
		vm:        vm,
		consensus: consensusMgr,
		wallet:    wallet,
		registry:  registry,
		ledger:    ledger,
		gas: map[string]uint64{
			"Stage77FailoverInit":      110,
			"Stage77FailoverRegister":  55,
			"Stage77FailoverHeartbeat": 45,
			"Stage77FailoverActive":    30,
			"Stage77FailoverReport":    75,
			"EnterpriseBootstrap":      120,
			"EnterpriseConsensusSync":  95,
			"EnterpriseWalletSeal":     60,
			"EnterpriseNodeAudit":      75,
			"EnterpriseAuthorityElect": 80,
		},
		failover: failover,
	}

	for _, opt := range opts {
		opt(orchestrator)
	}

	inserted, err := synn.EnsureGasSchedule(orchestrator.gas)
	if err != nil {
		return nil, fmt.Errorf("ensure gas schedule: %w", err)
	}
	orchestrator.pendingInsertions = inserted

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
	if o.failover != nil {
		o.failover.RegisterNode(FailoverNode{ID: node.Address, Role: role, Region: "authority"})
	}
	o.refreshDiagnostics(ctx)
	return node, nil
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

	diag := EnterpriseDiagnostics{
		Timestamp:         time.Now().UTC(),
		VMRunning:         o.vm.Status(),
		VMMode:            o.vm.Mode().String(),
		VMConcurrency:     o.vm.Concurrency(),
		ConsensusNetworks: len(o.consensus.ListNetworks()),
		AuthorityNodes:    len(o.registry.List()),
		WalletAddress:     o.wallet.Address,
		GasCoverage:       make(map[string]uint64, len(o.gas)),
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

	if o.failover != nil {
		report := o.failover.Report(ctx)
		diag.Resilience = &report
		if report.WalletAddress != "" {
			diag.WalletAddress = report.WalletAddress
		}
		if report.LedgerHeight > diag.LedgerHeight {
			diag.LedgerHeight = report.LedgerHeight
		}
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
