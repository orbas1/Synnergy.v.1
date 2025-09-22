package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	synn "synnergy"
	"synnergy/core"
	"synnergy/internal/security"
	tokens "synnergy/internal/tokens"
)

// runtime encapsulates the long-lived services required by the CLI and web
// consoles.  It wires core modules together so that commands can rely on a
// ready-to-use ledger, VM and consensus network without bespoke setup.
type runtime struct {
	ctx    context.Context
	cancel context.CancelFunc

	gasTable synn.GasTable
	vm       *core.SimpleVM
	ledger   *core.Ledger

	consensus *core.ConsensusNetworkManager
	wallet    *core.Wallet

	wg sync.WaitGroup
}

func bootstrapRuntime(parent context.Context) (*runtime, error) {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)

	if err := ensureGasCatalogue(); err != nil {
		cancel()
		return nil, err
	}

	vm := core.NewSimpleVM()
	if err := vm.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("start VM: %w", err)
	}

	ledger := core.NewLedger()
	consensus := core.NewConsensusNetworkManager()

	wallet, err := core.NewWallet()
	if err != nil {
		// Wallet creation failure should not abort startup, but it should
		// be visible.  Consumers can still retry via CLI commands.
		logrus.WithError(err).Warn("wallet bootstrap degraded")
	}

	rt := &runtime{
		ctx:       ctx,
		cancel:    cancel,
		gasTable:  synn.LoadGasTable(),
		vm:        vm,
		ledger:    ledger,
		consensus: consensus,
		wallet:    wallet,
	}

	if err := rt.preloadModules(); err != nil {
		cancel()
		return nil, err
	}

	rt.wg.Add(1)
	go func() {
		defer rt.wg.Done()
		<-ctx.Done()
		if err := vm.Stop(); err != nil {
			logrus.WithError(err).Warn("vm shutdown")
		}
	}()

	return rt, nil
}

func (rt *runtime) Shutdown() {
	if rt == nil {
		return
	}
	rt.cancel()
	rt.wg.Wait()
}

func (rt *runtime) preloadModules() error {
	// Shared registries
	_ = core.NewAuthorityNodeRegistry()
	_ = core.NewProposalManager()
	daoMgr := core.NewDAOManager()
	_ = core.NewDAOStaking(daoMgr)
	_ = core.NewDAOTokenLedger(daoMgr)

	// Banking and custodial flows use the shared ledger
	_ = core.NewBankInstitutionalNode("init", "init", rt.ledger)
	_ = core.NewCustodialNode("cli-custodian", "cli-custodian", rt.ledger)

	// Bridge and cross-chain connectivity
	_ = core.NewBridgeRegistry()
	_ = core.NewBridgeTransferManager()
	_ = core.NewChainConnectionManager()
	_ = core.NewProtocolRegistry()
	_ = core.NewCrossChainTxManager(rt.ledger)

	// Contract registry coupled to the bootstrap VM
	_ = core.NewContractRegistry(rt.vm)

	// Sandbox and security services
	_ = core.NewSandboxManager()
	_ = core.NewZeroTrustEngine()
	_ = security.NewSecretsManager()

	// Observability and specialised nodes
	_ = core.NewWarfareNode(core.NewNode("cli-war", "cli-war", rt.ledger))
	_ = core.NewWatchtowerNode("cli-watchtower", nil)
	_ = core.NewContentNetworkNode("cli-content", "cli")
	_ = core.NewNFTMarketplace()

	// Regulatory oversight
	_ = core.NewRegulatoryNode("cli-regnode", core.NewRegulatoryManager())

	// Token catalogue for CLI operations
	_ = tokens.NewSYN223Token("cli", "S223", "cli", 0)
	_ = tokens.NewSYN2700Token()
	_ = tokens.NewSYN3200Token(1)
	_ = tokens.NewSYN3600Token()
	_ = tokens.NewSYN3800Token(0)
	_ = tokens.NewSYN3900Token()
	_ = tokens.NewSYN500Token()
	_ = tokens.NewSYN5000Token()

	return nil
}

func ensureGasCatalogue() error {
	tbl := synn.LoadGasTable()
	var missing []string
	for _, name := range requiredGasOperations {
		cost, ok := tbl[name]
		if !ok {
			missing = append(missing, name)
			continue
		}
		if cost == 0 {
			missing = append(missing, name)
			continue
		}
		if err := synn.RegisterGasCost(name, cost); err != nil {
			return fmt.Errorf("register %s: %w", name, err)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("gas table missing entries: %s", strings.Join(missing, ", "))
	}
	return nil
}

var requiredGasOperations = []string{
	// Mining and governance
	"MineBlock",
	"CreateDAO",
	"UpdateMemberRole",
	"RenewAuthorityTerm",
	// Cross-chain bridges and protocols
	"RegisterBridge",
	"BridgeDeposit",
	"BridgeClaim",
	"OpenConnection",
	"CloseConnection",
	"LockMint",
	"BurnRelease",
	// Node and infrastructure operations
	"SetMode",
	"Stake",
	"Unstake",
	"Optimize",
	"SecureCommand",
	"TrackLogistics",
	"ShareTactical",
	"ReportFork",
	"Metrics",
	// Smart contract marketplace
	"DeploySmartContract",
	"TradeContract",
	"DeployTokenFaucetTemplate",
	"DeployStorageMarketTemplate",
	"DeployDAOGovernanceTemplate",
	"DeployNFTMintingTemplate",
	"DeployAIModelMarketTemplate",
	// Storage marketplace
	"CreateListing",
	"ListListings",
	"GetListing",
	"OpenDeal",
	"CloseDeal",
	"ListDeals",
	"GetDeal",
	"Storage_Pin",
	"Storage_Retrieve",
	"IPFS_Add",
	"IPFS_Get",
	"IPFS_Unpin",
	// NFT marketplace
	"MintNFT",
	"ListNFT",
	"BuyNFT",
	// Liquidity operations
	"Liquidity_Pool",
	"Liquidity_Pools",
	// Wallet integrations
	"NewWallet",
	"Sign",
	"VerifySignature",
	// Content management
	"RegisterContentNode",
	"UploadContent",
	"RetrieveContent",
	"ListContentNodes",
	// Monetary policy queries
	"BlockReward",
	"CirculatingSupply",
	"RemainingSupply",
	"InitialPrice",
	"AlphaFactor",
	"MinimumStake",
	// Kademlia and discovery
	"KademliaStore",
	"KademliaGet",
	"KademliaClosest",
	"KademliaDistance",
	// Stage 79 ledgers and governance enhancements
	"Ledger_StreamReplication",
	"Ledger_PrimaryElection",
	"Ledger_PrivacyEnvelope",
	"Consensus_StateDigest",
	"Consensus_AttestValidator",
	"Wallet_PrivacyEnvelope",
	"Node_Attest",
	"Authority_Attestation",
	"Loanpool_ComplianceDisburse",
	"CLI_CommandManifest",
}
