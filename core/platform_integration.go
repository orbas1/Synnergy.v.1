package core

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"
)

// PlatformIntegration coordinates core subsystems so higher level tooling can
// perform health checks and orchestrate workflows in a single pass. It wires
// together the virtual machine, consensus network manager, wallet, generic node
// and authority registry used by the enterprise CLI and web integrations.
type PlatformIntegration struct {
	vm        *SimpleVM
	node      *Node
	wallet    *Wallet
	authority *AuthorityNodeRegistry
	consensus *ConsensusNetworkManager
}

// IntegrationStatus captures the runtime state of each subsystem alongside the
// results of synthetic diagnostics executed during a health probe. The
// diagnostics map records human readable outcomes for quick inspection while
// Issues aggregates anything that requires operator attention.
type IntegrationStatus struct {
	Timestamp   time.Time           `json:"timestamp"`
	VM          VMStatus            `json:"vm"`
	Node        NodeStatus          `json:"node"`
	Wallet      WalletStatus        `json:"wallet"`
	Consensus   ConsensusStatus     `json:"consensus"`
	Authority   AuthorityStatus     `json:"authority"`
	Diagnostics map[string]string   `json:"diagnostics"`
	Issues      []string            `json:"issues"`
	Enterprise  EnterpriseReadiness `json:"enterprise"`
}

// DiagnosticCheck captures the outcome of an enterprise assurance probe.
type DiagnosticCheck struct {
	Healthy bool   `json:"healthy"`
	Detail  string `json:"detail"`
	Latency string `json:"latency"`
}

// EnterpriseReadiness aggregates enterprise grade assurances required by the
// integration workflow. Each dimension maps to a core platform capability.
type EnterpriseReadiness struct {
	Security         DiagnosticCheck `json:"security"`
	Scalability      DiagnosticCheck `json:"scalability"`
	Privacy          DiagnosticCheck `json:"privacy"`
	Governance       DiagnosticCheck `json:"governance"`
	Interoperability DiagnosticCheck `json:"interoperability"`
	Compliance       DiagnosticCheck `json:"compliance"`
}

// VMStatus reports runtime and capacity details for the virtual machine.
type VMStatus struct {
	Running     bool   `json:"running"`
	Mode        string `json:"mode"`
	Concurrency int    `json:"concurrency"`
}

// NodeStatus summarises the blockchain node that backs integration tests.
type NodeStatus struct {
	ID            string `json:"id"`
	BlockHeight   int    `json:"block_height"`
	PendingTx     int    `json:"pending_transactions"`
	LastBlockHash string `json:"last_block_hash"`
}

// WalletStatus reports the wallet used for signing integration transactions.
type WalletStatus struct {
	Address string `json:"address"`
}

// ConsensusStatus captures the cross-consensus manager state exposed to the CLI.
type ConsensusStatus struct {
	Weights          ConsensusWeights `json:"weights"`
	Networks         int              `json:"networks"`
	AuthorizedRelays int              `json:"authorized_relays"`
}

// AuthorityStatus records the number of authority nodes that are currently wired
// into the orchestrated environment.
type AuthorityStatus struct {
	Registered int `json:"registered"`
}

// NewPlatformIntegration constructs a ready-to-use integration surface. It starts
// a heavy virtual machine profile for stress testing, provisions a wallet used
// for synthetic transactions, allocates a ledger-backed node and authorizes the
// wallet as a relayer for cross-consensus operations.
func NewPlatformIntegration() (*PlatformIntegration, error) {
	vm := NewSimpleVM(VMHeavy)
	if err := vm.Start(); err != nil {
		return nil, fmt.Errorf("start vm: %w", err)
	}

	wallet, err := NewWallet()
	if err != nil {
		return nil, fmt.Errorf("wallet: %w", err)
	}

	ledger := NewLedger()
	node := NewNode("enterprise-node", "enterprise", ledger)
	authority := NewAuthorityNodeRegistry()
	if _, err := authority.Register(wallet.Address, "integration-orchestrator"); err != nil {
		return nil, fmt.Errorf("authority register: %w", err)
	}

	consensus := NewConsensusNetworkManager()
	consensus.AuthorizeRelayer(wallet.Address)

	return &PlatformIntegration{
		vm:        vm,
		node:      node,
		wallet:    wallet,
		authority: authority,
		consensus: consensus,
	}, nil
}

// Close stops background resources associated with the integration surface.
func (p *PlatformIntegration) Close() error {
	if p == nil || p.vm == nil {
		return nil
	}
	return p.vm.Stop()
}

// Diagnostics executes a battery of synthetic transactions to ensure the
// subsystems remain tightly integrated. It returns a status snapshot that can be
// serialized for dashboards or CLI tooling.
func (p *PlatformIntegration) Diagnostics(ctx context.Context) IntegrationStatus {
	if ctx == nil {
		ctx = context.Background()
	}
	start := time.Now()
	status := IntegrationStatus{
		Timestamp:   time.Now().UTC(),
		Diagnostics: make(map[string]string),
		VM: VMStatus{
			Running:     p.vm.Status(),
			Mode:        p.vm.ModeString(),
			Concurrency: p.vm.Concurrency(),
		},
		Node: NodeStatus{
			ID:          p.node.ID,
			BlockHeight: p.node.BlockHeight(),
			PendingTx:   p.node.PendingTransactions(),
		},
		Wallet: WalletStatus{
			Address: p.wallet.Address,
		},
		Consensus: ConsensusStatus{
			Weights:          p.node.Consensus.Weights,
			Networks:         len(p.consensus.ListNetworks()),
			AuthorizedRelays: len(p.consensus.AuthorizedRelayers()),
		},
		Authority: AuthorityStatus{
			Registered: len(p.authority.List()),
		},
	}

	recordContext := func(stage string) bool {
		if err := ctx.Err(); err != nil {
			issue := fmt.Sprintf("%s cancelled: %v", stage, err)
			status.Diagnostics[stage] = issue
			status.Issues = append(status.Issues, issue)
			return false
		}
		return true
	}

	if status.VM.Running {
		status.Diagnostics["vm"] = "runtime online"
	} else {
		issue := "virtual machine offline"
		status.Diagnostics["vm"] = issue
		status.Issues = append(status.Issues, issue)
	}

	ledger := p.node.Ledger
	ledger.Credit(p.wallet.Address, 1_000_000)
	if err := p.node.SetStake(p.wallet.Address, MinStake); err != nil {
		issue := fmt.Sprintf("validator stake: %v", err)
		status.Diagnostics["validator"] = issue
		status.Issues = append(status.Issues, issue)
	} else {
		status.Diagnostics["validator"] = "staking ready"
	}

	if !recordContext("wallet_generation") {
		status.Diagnostics["overall"] = "integration cancelled"
		status.Diagnostics["latency"] = time.Since(start).String()
		return status
	}
	targetWallet, err := NewWallet()
	if err != nil {
		issue := fmt.Sprintf("wallet generation: %v", err)
		status.Diagnostics["wallet"] = issue
		status.Issues = append(status.Issues, issue)
	} else {
		tx := NewTransaction(p.wallet.Address, targetWallet.Address, 500, 10, 0)
		if err := p.node.AddTransaction(tx); err != nil {
			issue := fmt.Sprintf("transaction enqueue: %v", err)
			status.Diagnostics["transaction"] = issue
			status.Issues = append(status.Issues, issue)
		} else {
			block := p.node.MineBlock()
			if block == nil {
				issue := "consensus failed to mine diagnostic block"
				status.Diagnostics["transaction"] = issue
				status.Issues = append(status.Issues, issue)
			} else {
				status.Diagnostics["transaction"] = fmt.Sprintf("mined block %s", block.Hash)
				status.Node.BlockHeight = p.node.BlockHeight()
				status.Node.PendingTx = p.node.PendingTransactions()
				status.Node.LastBlockHash = block.Hash
				if h, ok := p.node.Ledger.GetBlock(status.Node.BlockHeight); ok {
					status.Node.LastBlockHash = h.Hash
				}
			}
		}
	}

	if !recordContext("consensus") {
		status.Diagnostics["overall"] = "integration cancelled"
		status.Diagnostics["latency"] = time.Since(start).String()
		return status
	}
	if _, err := p.consensus.RegisterNetwork("synnergy-main", "synnergy-enterprise", p.wallet.Address); err != nil {
		issue := fmt.Sprintf("consensus registration: %v", err)
		status.Diagnostics["consensus"] = issue
		status.Issues = append(status.Issues, issue)
	} else {
		status.Consensus.Networks = len(p.consensus.ListNetworks())
		status.Diagnostics["consensus"] = "enterprise relay active"
	}

	prefix := p.wallet.Address
	if len(prefix) > 8 {
		prefix = prefix[:8]
	}
	authorityID := fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	registeredAuthority := ""
	if !recordContext("authority") {
		status.Diagnostics["overall"] = "integration cancelled"
		status.Diagnostics["latency"] = time.Since(start).String()
		return status
	}
	if _, err := p.authority.Register(authorityID, "observer"); err != nil {
		issue := fmt.Sprintf("authority registration: %v", err)
		status.Diagnostics["authority"] = issue
		status.Issues = append(status.Issues, issue)
	} else {
		status.Authority.Registered = len(p.authority.List())
		status.Diagnostics["authority"] = "observer enrolled"
		registeredAuthority = authorityID
	}

	status.Enterprise.Security = p.runEnterpriseCheck(ctx, &status, "enterprise_security", p.securityCheck)
	status.Enterprise.Scalability = p.runEnterpriseCheck(ctx, &status, "enterprise_scalability", p.scalabilityCheck)
	status.Node.BlockHeight = p.node.BlockHeight()
	status.Node.PendingTx = p.node.PendingTransactions()
	if h, ok := p.node.Ledger.GetBlock(status.Node.BlockHeight); ok {
		status.Node.LastBlockHash = h.Hash
	}
	status.Enterprise.Privacy = p.runEnterpriseCheck(ctx, &status, "enterprise_privacy", p.privacyCheck)
	status.Enterprise.Governance = p.runEnterpriseCheck(ctx, &status, "enterprise_governance", func(checkCtx context.Context) (string, error) {
		return p.governanceCheck(checkCtx, registeredAuthority)
	})
	status.Enterprise.Interoperability = p.runEnterpriseCheck(ctx, &status, "enterprise_interoperability", p.interoperabilityCheck)
	status.Enterprise.Compliance = p.runEnterpriseCheck(ctx, &status, "enterprise_compliance", p.complianceCheck)

	if len(status.Issues) == 0 {
		status.Diagnostics["overall"] = "integration healthy"
	} else {
		status.Diagnostics["overall"] = "integration requires attention"
	}
	status.Diagnostics["latency"] = time.Since(start).String()

	return status
}

func (p *PlatformIntegration) runEnterpriseCheck(ctx context.Context, status *IntegrationStatus, key string, fn func(context.Context) (string, error)) DiagnosticCheck {
	checkStart := time.Now()
	detail, err := fn(ctx)
	if err != nil {
		detail = err.Error()
		status.Issues = append(status.Issues, fmt.Sprintf("%s: %s", key, detail))
	}
	if status.Diagnostics == nil {
		status.Diagnostics = make(map[string]string)
	}
	status.Diagnostics[key] = detail
	return DiagnosticCheck{
		Healthy: err == nil,
		Detail:  detail,
		Latency: time.Since(checkStart).String(),
	}
}

func (p *PlatformIntegration) securityCheck(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	tx := NewTransaction(p.wallet.Address, p.wallet.Address, 1, 1, 42)
	if _, err := p.wallet.Sign(tx); err != nil {
		return "", fmt.Errorf("sign transaction: %w", err)
	}
	if !tx.Verify(&p.wallet.PrivateKey.PublicKey) {
		return "", errors.New("signature verification failed")
	}
	tmp, err := os.CreateTemp("", "syn-integration-wallet-*.json")
	if err != nil {
		return "", fmt.Errorf("temp wallet: %w", err)
	}
	defer os.Remove(tmp.Name())
	if err := p.wallet.Save(tmp.Name(), "enterprise-grade"); err != nil {
		return "", fmt.Errorf("save wallet: %w", err)
	}
	loaded, err := LoadWallet(tmp.Name(), "enterprise-grade")
	if err != nil {
		return "", fmt.Errorf("load wallet: %w", err)
	}
	if loaded.Address != p.wallet.Address {
		return "", errors.New("wallet address mismatch after reload")
	}
	return "digital signatures verified and wallet storage encrypted", nil
}

func (p *PlatformIntegration) scalabilityCheck(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	const burst = 25
	burstStart := time.Now()
	for i := 0; i < burst; i++ {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		target := fmt.Sprintf("throughput-%d", i)
		tx := NewTransaction(p.wallet.Address, target, 10, 1, uint64(i+1))
		if err := p.node.AddTransaction(tx); err != nil {
			return "", fmt.Errorf("queue tx %d: %w", i, err)
		}
	}
	block := p.node.MineBlock()
	if block == nil {
		return "", errors.New("throughput block not mined")
	}
	processed := 0
	for _, sb := range block.SubBlocks {
		processed += len(sb.Transactions)
	}
	elapsed := time.Since(burstStart)
	if elapsed <= 0 {
		elapsed = time.Millisecond
	}
	rate := float64(processed) / elapsed.Seconds()
	return fmt.Sprintf("processed %d tx @ %.2f tx/s", processed, rate), nil
}

func (p *PlatformIntegration) privacyCheck(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if !p.consensus.IsRelayerAuthorized(p.wallet.Address) {
		return "", errors.New("integration wallet not authorized as relayer")
	}
	other, err := NewWallet()
	if err != nil {
		return "", fmt.Errorf("wallet provision: %w", err)
	}
	if _, err := p.consensus.RegisterNetwork("unauthorised", "target", other.Address); err == nil {
		return "", errors.New("unauthorized relayer allowed to register network")
	}
	return "consensus relayer whitelist enforced", nil
}

func (p *PlatformIntegration) governanceCheck(ctx context.Context, candidate string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if candidate == "" {
		return "", errors.New("no authority candidate registered")
	}
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", fmt.Errorf("generate voter keys: %w", err)
	}
	voter := hex.EncodeToString(pub)
	sig := ed25519.Sign(priv, []byte(candidate))
	if err := p.authority.Vote(voter, candidate, sig, pub); err != nil {
		return "", fmt.Errorf("record vote: %w", err)
	}
	p.authority.RemoveVote(voter, candidate)
	return "governance vote cryptography validated", nil
}

func (p *PlatformIntegration) interoperabilityCheck(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	mgr := NewCrossChainTxManager(p.node.Ledger)
	mgr.AuthorizeRelayer(p.wallet.Address)
	if _, err := mgr.LockMint(1, p.wallet.Address, "interop-destination", "SYN", 50, "proof"); err != nil {
		return "", fmt.Errorf("lock mint: %w", err)
	}
	if _, err := mgr.BurnRelease(1, p.wallet.Address, p.wallet.Address, "SYN", 25); err != nil {
		return "", fmt.Errorf("burn release: %w", err)
	}
	transfers := mgr.ListTransfers()
	return fmt.Sprintf("cross-chain relay active (%d records)", len(transfers)), nil
}

func (p *PlatformIntegration) complianceCheck(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	mgr := NewRegulatoryManager()
	if err := mgr.AddRegulation(Regulation{ID: "aml_limit", Jurisdiction: "global", Description: "transaction cap", MaxAmount: 2000}); err != nil {
		return "", fmt.Errorf("add regulation: %w", err)
	}
	node := NewRegulatoryNode("integration-regulator", mgr)
	node.RegisterWallet(p.wallet)
	tx := NewTransaction(p.wallet.Address, "enterprise-beneficiary", 500, 1, 900)
	if _, err := p.wallet.Sign(tx); err != nil {
		return "", fmt.Errorf("sign compliant tx: %w", err)
	}
	if err := node.ApproveTransaction(*tx); err != nil {
		return "", fmt.Errorf("expected approval: %w", err)
	}
	highValue := NewTransaction(p.wallet.Address, "enterprise-beneficiary", 5000, 1, 901)
	if _, err := p.wallet.Sign(highValue); err != nil {
		return "", fmt.Errorf("sign violation tx: %w", err)
	}
	if err := node.ApproveTransaction(*highValue); err == nil {
		return "", errors.New("regulation violation not detected")
	}
	logs := node.Logs(p.wallet.Address)
	if len(logs) == 0 {
		return "", errors.New("regulatory violation not logged")
	}
	return fmt.Sprintf("compliance guard flagged %d issue(s)", len(logs)), nil
}
