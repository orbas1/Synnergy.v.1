package core

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"time"

	synn "synnergy"
	"synnergy/internal/telemetry"
)

// EnterpriseOption configures the enterprise orchestrator at construction time.
type EnterpriseOption func(*EnterpriseOrchestrator)

// EnterpriseBootstrapConfig captures orchestration preferences for Stage 79
// enterprise bootstrap flows. Operators can declare consensus and governance
// profiles, request ledger replication and provide additional authority nodes
// that must be registered atomically with the orchestrator wallet.
type EnterpriseBootstrapConfig struct {
	NodeID            string
	Address           string
	ConsensusProfile  string
	GovernanceProfile string
	EnableReplication bool
	EnableRegulator   bool
	Authorities       map[string]string
}

// EnterpriseBootstrapResult summarises the outcome of a bootstrap invocation.
// It includes diagnostics so CLI, web and automation clients can surface a
// consistent view of ledger height, consensus reachability and authority
// membership immediately after orchestration.
type EnterpriseBootstrapResult struct {
	NodeID             string                `json:"nodeId"`
	Address            string                `json:"address"`
	ConsensusNetworkID int                   `json:"consensusNetworkId"`
	AuthorityNodes     []string              `json:"authorityNodes"`
	ReplicationEnabled bool                  `json:"replicationEnabled"`
	WalletAddress      string                `json:"walletAddress"`
	BootstrapSignature string                `json:"bootstrapSignature"`
	Diagnostics        EnterpriseDiagnostics `json:"diagnostics"`
}

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
	BootstrapNodes    int               `json:"bootstrapNodes"`
	ReplicationActive bool              `json:"replicationActive"`
	GasCoverage       map[string]uint64 `json:"gasCoverage"`
	MissingOpcodes    []string          `json:"missingOpcodes"`
	InsertedOpcodes   []string          `json:"insertedOpcodes,omitempty"`
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
	replicators       []*Replicator
	nodes             map[string]*Node
	gas               map[string]uint64
	last              EnterpriseDiagnostics
	pendingInsertions []string
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
		vm:          vm,
		consensus:   NewConsensusNetworkManager(),
		wallet:      wallet,
		registry:    NewAuthorityNodeRegistry(),
		ledger:      NewLedger(),
		replicators: []*Replicator{},
		nodes:       make(map[string]*Node),
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

// BootstrapNetwork provisions a new enterprise-ready node, wiring it into the
// consensus mesh, registering authorities and enabling ledger replication in a
// single transaction. The orchestrator wallet signs the bootstrap to provide a
// verifiable audit trail for governance and regulatory systems.
func (o *EnterpriseOrchestrator) BootstrapNetwork(ctx context.Context, cfg EnterpriseBootstrapConfig) (EnterpriseBootstrapResult, error) {
	if cfg.NodeID == "" {
		return EnterpriseBootstrapResult{}, fmt.Errorf("bootstrap: node id required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if cfg.Address == "" {
		cfg.Address = cfg.NodeID
	}

	_, span := telemetry.Tracer("core.enterprise").Start(ctx, "EnterpriseOrchestrator.BootstrapNetwork")
	defer span.End()

	if !o.vm.Status() {
		if err := o.vm.Start(); err != nil {
			span.RecordError(err)
			return EnterpriseBootstrapResult{}, fmt.Errorf("start vm: %w", err)
		}
	}

	node := NewNode(cfg.NodeID, cfg.Address, o.ledger)
	if cfg.EnableRegulator {
		regulator := NewRegulatoryNode("bootstrap-regulator", NewRegulatoryManager())
		node.Consensus.SetRegulatoryNode(regulator)
	}

	var replicator *Replicator
	if cfg.EnableReplication {
		replicator = NewReplicator(o.ledger)
		replicator.Start()
		if _, hash := o.ledger.Head(); hash != "" {
			replicator.ReplicateBlock(hash)
		} else {
			replicator.ReplicateBlock("genesis")
		}
		o.mu.Lock()
		o.replicators = append(o.replicators, replicator)
		o.mu.Unlock()
	}

	source := cfg.ConsensusProfile
	if source == "" {
		source = "Synnergy-PBFT"
	}
	target := cfg.GovernanceProfile
	if target == "" {
		target = "Synnergy-PBFT"
	}
	consensusID, err := o.consensus.RegisterNetwork(source, target, o.wallet.Address)
	if err != nil {
		span.RecordError(err)
		return EnterpriseBootstrapResult{}, err
	}

	authorityMap := make(map[string]string, len(cfg.Authorities)+1)
	authorityMap[o.wallet.Address] = "orchestrator"
	for addr, role := range cfg.Authorities {
		if addr == "" {
			continue
		}
		authorityMap[addr] = role
	}

	registered := make([]string, 0, len(authorityMap))
	for addr, role := range authorityMap {
		if addr == o.wallet.Address {
			registered = append(registered, addr)
			continue
		}
		if !o.registry.IsAuthorityNode(addr) {
			if _, err := o.registry.Register(addr, role); err != nil {
				span.RecordError(err)
				return EnterpriseBootstrapResult{}, err
			}
		}
		registered = append(registered, addr)
	}
	sort.Strings(registered)

	o.mu.Lock()
	o.nodes[cfg.NodeID] = node
	o.mu.Unlock()

	signature, err := o.wallet.SignMessage([]byte("bootstrap:" + cfg.NodeID + ":" + cfg.Address))
	if err != nil {
		span.RecordError(err)
		return EnterpriseBootstrapResult{}, fmt.Errorf("sign bootstrap: %w", err)
	}

	diag := o.refreshDiagnostics(ctx)
	result := EnterpriseBootstrapResult{
		NodeID:             cfg.NodeID,
		Address:            cfg.Address,
		ConsensusNetworkID: consensusID,
		AuthorityNodes:     registered,
		ReplicationEnabled: cfg.EnableReplication && replicator != nil,
		WalletAddress:      o.wallet.Address,
		BootstrapSignature: hex.EncodeToString(signature),
		Diagnostics:        diag,
	}
	return result, nil
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
	o.mu.RLock()
	diag.BootstrapNodes = len(o.nodes)
	replicators := make([]*Replicator, len(o.replicators))
	copy(replicators, o.replicators)
	o.mu.RUnlock()
	for _, rep := range replicators {
		if rep != nil && rep.Status() {
			diag.ReplicationActive = true
			break
		}
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

	o.mu.Lock()
	if len(o.pendingInsertions) > 0 {
		diag.InsertedOpcodes = append(diag.InsertedOpcodes, o.pendingInsertions...)
		o.pendingInsertions = nil
	}
	o.last = diag
	o.mu.Unlock()

	return diag
}
