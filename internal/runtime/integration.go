package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	synn "synnergy"
	"synnergy/core"
)

// IntegrationConfig captures runtime wiring details required to expose the core
// modules consistently across the CLI, virtual machine and web surfaces.
type IntegrationConfig struct {
	NodeID              string
	NodeAddress         string
	AuthorityAddress    string
	AuthorityDepartment string
	GasLimit            uint64
	CriticalOpcodes     []string
	MonitorInterval     time.Duration
}

// RuntimeIntegration orchestrates consensus, wallet, node and VM components so
// that higher level entrypoints (CLI, APIs and the web interface) interact with
// a cohesive stack instead of ad-hoc singletons.
type RuntimeIntegration struct {
	ledger          *core.Ledger
	node            *core.Node
	vm              *synn.SimpleVM
	consensus       *core.SynnergyConsensus
	wallet          *core.Wallet
	authority       *core.GovernmentAuthorityNode
	regNode         *core.RegulatoryNode
	gasLimit        uint64
	monitorInterval time.Duration
	criticalOpcodes []string

	mu       sync.RWMutex
	started  bool
	cancel   context.CancelFunc
	healthCh chan RuntimeHealth
}

// RuntimeHealth summarises the current state of the integrated stack. It is
// emitted on each monitoring cycle and may be consumed by dashboards or tests.
type RuntimeHealth struct {
	NodeID        string
	Started       bool
	VMRunning     bool
	LedgerHeight  int
	LastBlockHash string
	WalletAddress string
	ConsensusPoW  float64
	ConsensusPoS  float64
	ConsensusPoH  float64
	Timestamp     time.Time
}

// NewRuntimeIntegration wires together the ledger, node, consensus engine,
// wallet, authority node and regulatory node. It validates configuration and
// pre-registers the operator wallet with the regulatory node so transactions can
// be approved end-to-end.
func NewRuntimeIntegration(cfg IntegrationConfig) (*RuntimeIntegration, error) {
	if cfg.NodeID == "" {
		return nil, errors.New("node id required")
	}
	if cfg.NodeAddress == "" {
		return nil, errors.New("node address required")
	}
	if cfg.AuthorityAddress == "" {
		return nil, errors.New("authority address required")
	}

	ledger := core.NewLedger()
	node := core.NewNode(cfg.NodeID, cfg.NodeAddress, ledger)
	vm := synn.NewSimpleVM()
	consensus := node.Consensus
	regMgr := core.NewRegulatoryManager()
	regNode := core.NewRegulatoryNode(cfg.NodeID+"-reg", regMgr)
	consensus.SetRegulatoryNode(regNode)

	wallet, err := core.NewWallet()
	if err != nil {
		return nil, fmt.Errorf("wallet init: %w", err)
	}
	regNode.RegisterWallet(wallet)
	if err := node.SetStake(wallet.Address, core.MinStake); err != nil {
		return nil, fmt.Errorf("stake bootstrap: %w", err)
	}

	authority := core.NewGovernmentAuthorityNode(cfg.AuthorityAddress, "governor", cfg.AuthorityDepartment)

	gasLimit := cfg.GasLimit
	if gasLimit == 0 {
		gasLimit = 10_000
	}
	interval := cfg.MonitorInterval
	if interval == 0 {
		interval = 5 * time.Second
	}

	ri := &RuntimeIntegration{
		ledger:          ledger,
		node:            node,
		vm:              vm,
		consensus:       consensus,
		wallet:          wallet,
		authority:       authority,
		regNode:         regNode,
		gasLimit:        gasLimit,
		monitorInterval: interval,
		criticalOpcodes: append([]string(nil), cfg.CriticalOpcodes...),
		healthCh:        make(chan RuntimeHealth, 1),
	}
	return ri, nil
}

// Ledger exposes the shared ledger for callers that need to hydrate balances or
// inspect block history.
func (ri *RuntimeIntegration) Ledger() *core.Ledger { return ri.ledger }

// Node returns the fully configured node instance.
func (ri *RuntimeIntegration) Node() *core.Node { return ri.node }

// Wallet exposes the operator wallet.
func (ri *RuntimeIntegration) Wallet() *core.Wallet { return ri.wallet }

// WalletAddress returns the public address bound to the runtime wallet.
func (ri *RuntimeIntegration) WalletAddress() string { return ri.wallet.Address }

// HealthChannel returns a buffered channel that receives health snapshots each
// monitoring cycle. Consumers may safely range over the channel until Stop is
// invoked.
func (ri *RuntimeIntegration) HealthChannel() <-chan RuntimeHealth { return ri.healthCh }

// EnsureGasSchedule validates that every critical opcode has a gas price and
// registers the cached price inside the long-lived gas table used by CLI and VM
// helpers.
func (ri *RuntimeIntegration) EnsureGasSchedule() error {
	if err := synn.EnsureGasCosts(ri.criticalOpcodes); err != nil {
		return err
	}
	for _, name := range ri.criticalOpcodes {
		if err := synn.RegisterGasCost(name, synn.GasCost(name)); err != nil {
			return fmt.Errorf("register %s: %w", name, err)
		}
	}
	return nil
}

// Start launches the virtual machine and monitoring goroutine. It is safe to
// call multiple times.
func (ri *RuntimeIntegration) Start(ctx context.Context) error {
	ri.mu.Lock()
	if ri.started {
		ri.mu.Unlock()
		return nil
	}
	if err := ri.vm.Start(); err != nil {
		ri.mu.Unlock()
		return err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	cctx, cancel := context.WithCancel(ctx)
	ri.cancel = cancel
	ri.started = true
	ri.mu.Unlock()

	go ri.monitor(cctx)
	return nil
}

// Stop halts monitoring and the underlying VM.
func (ri *RuntimeIntegration) Stop() {
	ri.mu.Lock()
	if !ri.started {
		ri.mu.Unlock()
		return
	}
	ri.started = false
	if ri.cancel != nil {
		ri.cancel()
		ri.cancel = nil
	}
	ri.mu.Unlock()
	_ = ri.vm.Stop()
}

func (ri *RuntimeIntegration) monitor(ctx context.Context) {
	ticker := time.NewTicker(ri.monitorInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ri.emitHealth()
		}
	}
}

func (ri *RuntimeIntegration) emitHealth() {
	height, hash := ri.ledger.Head()
	snapshot := RuntimeHealth{
		NodeID:        ri.node.ID,
		Started:       ri.started,
		VMRunning:     ri.vm.Status(),
		LedgerHeight:  height,
		LastBlockHash: hash,
		WalletAddress: ri.wallet.Address,
		ConsensusPoW:  ri.consensus.Weights.PoW,
		ConsensusPoS:  ri.consensus.Weights.PoS,
		ConsensusPoH:  ri.consensus.Weights.PoH,
		Timestamp:     time.Now().UTC(),
	}
	select {
	case ri.healthCh <- snapshot:
	default:
		// Keep the latest snapshot without blocking producers.
		select {
		case <-ri.healthCh:
		default:
		}
		ri.healthCh <- snapshot
	}
}

// SubmitTransaction signs the transaction with the orchestrated wallet, submits
// it for regulatory approval and enqueues it in the node's mempool once
// validated.
func (ri *RuntimeIntegration) SubmitTransaction(ctx context.Context, tx *core.Transaction) error {
	if tx == nil {
		return errors.New("transaction required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	ri.mu.RLock()
	started := ri.started
	ri.mu.RUnlock()
	if !started {
		return errors.New("runtime integration not started")
	}
	if tx.From == "" {
		tx.From = ri.wallet.Address
	}
	if _, err := ri.wallet.Sign(tx); err != nil {
		return fmt.Errorf("sign transaction: %w", err)
	}
	if err := ri.regNode.ApproveTransaction(*tx); err != nil {
		return fmt.Errorf("regulatory approval failed: %w", err)
	}
	return ri.node.AddTransaction(tx)
}

// ExecuteProgram routes bytecode through the shared virtual machine so that CLI
// and web flows can evaluate contracts consistently. It enforces a shared gas
// limit derived from the integration configuration.
func (ri *RuntimeIntegration) ExecuteProgram(ctx context.Context, wasm []byte, method string, args []byte) ([]byte, uint64, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	ri.mu.RLock()
	started := ri.started
	ri.mu.RUnlock()
	if !started {
		return nil, 0, errors.New("runtime integration not started")
	}
	return ri.vm.ExecuteContext(ctx, wasm, method, args, ri.gasLimit)
}

// HealthReport provides the most recent snapshot without waiting on the channel.
func (ri *RuntimeIntegration) HealthReport() RuntimeHealth {
	select {
	case h := <-ri.healthCh:
		// Requeue the snapshot for future consumers.
		select {
		case ri.healthCh <- h:
		default:
		}
		return h
	default:
		height, hash := ri.ledger.Head()
		return RuntimeHealth{
			NodeID:        ri.node.ID,
			Started:       ri.started,
			VMRunning:     ri.vm.Status(),
			LedgerHeight:  height,
			LastBlockHash: hash,
			WalletAddress: ri.wallet.Address,
			ConsensusPoW:  ri.consensus.Weights.PoW,
			ConsensusPoS:  ri.consensus.Weights.PoS,
			ConsensusPoH:  ri.consensus.Weights.PoH,
			Timestamp:     time.Now().UTC(),
		}
	}
}

// CriticalOpcodes exposes the configured opcode catalogue.
func (ri *RuntimeIntegration) CriticalOpcodes() []string {
	return append([]string(nil), ri.criticalOpcodes...)
}
