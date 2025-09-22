package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	synn "synnergy"
	"synnergy/cli"
	"synnergy/core"
	"synnergy/internal/config"
)

func main() {
	gotenv.Load()
	otel.SetTracerProvider(trace.NewNoopTracerProvider())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfgPath := os.Getenv("SYN_CONFIG")
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := configureLogging(cfg); err != nil {
		logrus.Fatal(err)
	}

	rt, err := bootstrapRuntime(ctx)
	if err != nil {
		logrus.Fatalf("bootstrap runtime: %v", err)
	}
	defer rt.Shutdown()

	rt.vm.RegisterHook(func(trace core.ExecutionTrace) {
		if trace.Err != nil {
			logrus.WithError(trace.Err).WithField("opcode", trace.Name).Warn("vm execution error")
		}
	})

	if err := registerEnterpriseGasMetadata(); err != nil {
		logrus.Fatalf("gas metadata: %v", err)
	}

	orch, err := core.NewEnterpriseOrchestrator(ctx)
	if err != nil {
		logrus.Fatalf("enterprise orchestrator: %v", err)
	}
	diag, err := orch.EnterpriseBootstrap(ctx)
	if err != nil {
		logrus.Fatalf("enterprise bootstrap: %v", err)
	}
	cli.InjectEnterpriseOrchestrator(orch)

	logrus.WithFields(logrus.Fields{
		"wallet":             diag.WalletAddress,
		"sealed":             diag.WalletSealed,
		"consensus_networks": diag.ConsensusNetworks,
		"authority_nodes":    diag.AuthorityNodes,
		"gas_synced_at":      diag.GasLastSyncedAt,
	}).Info("enterprise orchestrator initialised")

	logrus.Infof("starting Synnergy in %s mode on %s:%d", cfg.Environment, cfg.Server.Host, cfg.Server.Port)

	if err := cli.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func configureLogging(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}
	lvl, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	logrus.SetLevel(lvl)
	logrus.SetReportCaller(cfg.Log.IncludeCaller)

	switch strings.ToLower(cfg.Log.Format) {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	default:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	outputs := make([]io.Writer, 0, len(cfg.Log.Outputs))
	for _, out := range cfg.Log.Outputs {
		switch strings.ToLower(out) {
		case "stderr":
			outputs = append(outputs, os.Stderr)
		case "stdout":
			outputs = append(outputs, os.Stdout)
		default:
			outputs = append(outputs, os.Stdout)
		}
	}
	if len(outputs) == 0 {
		outputs = []io.Writer{os.Stdout}
	}
	logrus.SetOutput(io.MultiWriter(outputs...))
	return nil
}

func registerEnterpriseGasMetadata() error {
	synn.LoadGasTable()
	stage78Gas := map[string]uint64{
		"EnterpriseBootstrap":      120,
		"EnterpriseConsensusSync":  95,
		"EnterpriseWalletSeal":     60,
		"EnterpriseNodeAudit":      75,
		"EnterpriseAuthorityElect": 80,
	}
	if inserted, err := synn.EnsureGasSchedule(stage78Gas); err != nil {
		return fmt.Errorf("stage 78 gas sync failed: %w", err)
	} else if len(inserted) > 0 {
		logrus.Infof("registered %d stage 78 opcodes", len(inserted))
	}


	categories := []struct {
		category    string
		description string
		names       []string
	}{
		{"consensus", "Core consensus lifecycle operations", []string{"MineBlock"}},
		{"dao", "DAO creation and authority renewal", []string{"CreateDAO", "UpdateMemberRole", "RenewAuthorityTerm"}},
		{"cross-chain", "Stage 24 cross-chain operations", []string{"RegisterBridge", "BridgeDeposit", "BridgeClaim", "OpenConnection", "CloseConnection", "LockMint", "BurnRelease"}},
		{"node", "Stage 25 node and infrastructure operations", []string{"SetMode", "Stake", "Unstake", "Optimize", "SecureCommand", "TrackLogistics", "ShareTactical", "ReportFork", "Metrics"}},
		{"templates", "Stage 29 contract templates", []string{"DeployTokenFaucetTemplate", "DeployStorageMarketTemplate", "DeployDAOGovernanceTemplate", "DeployNFTMintingTemplate", "DeployAIModelMarketTemplate"}},
		{"marketplace", "Stage 34 marketplace settlement", []string{"DeploySmartContract", "TradeContract"}},
		{"storage", "Stage 35 storage marketplace operations", []string{"CreateListing", "ListListings", "GetListing", "OpenDeal", "CloseDeal", "ListDeals", "GetDeal", "Storage_Pin", "Storage_Retrieve", "IPFS_Add", "IPFS_Get", "IPFS_Unpin"}},
		{"nft", "Stage 36 NFT marketplace operations", []string{"MintNFT", "ListNFT", "BuyNFT"}},
		{"dex", "Stage 39 liquidity view operations", []string{"Liquidity_Pool", "Liquidity_Pools"}},
		{"wallet", "Wallet lifecycle operations", []string{"NewWallet", "Sign", "VerifySignature"}},
		{"content", "Stage 59 content registry operations", []string{"RegisterContentNode", "UploadContent", "RetrieveContent", "ListContentNodes"}},
		{"monetary", "Stage 40 monetary policy queries", []string{"BlockReward", "CirculatingSupply", "RemainingSupply", "InitialPrice", "AlphaFactor", "MinimumStake"}},
		{"p2p", "Stage 67 Kademlia routing operations", []string{"KademliaStore", "KademliaGet", "KademliaClosest", "KademliaDistance"}},
		{"orchestrator", "Stage 78 enterprise orchestrator operations", []string{"EnterpriseBootstrap", "EnterpriseConsensusSync", "EnterpriseWalletSeal", "EnterpriseNodeAudit", "EnterpriseAuthorityElect"}},
	}

	for _, entry := range categories {
		for _, name := range entry.names {
=======
	enterpriseSpecialGas := map[string]uint64{
		"EnterpriseSpecialAttach":    110,
		"EnterpriseSpecialDetach":    55,
		"EnterpriseSpecialBroadcast": 145,
		"EnterpriseSpecialSnapshot":  40,
		"EnterpriseSpecialLedger":    30,
	}
	if inserted, err := synn.EnsureGasSchedule(enterpriseSpecialGas); err != nil {
		logrus.Fatalf("enterprise special gas sync failed: %v", err)
	} else if len(inserted) > 0 {
		logrus.Infof("registered %d enterprise special opcodes", len(inserted))
	}
	register := func(category, description string, names ...string) {
		for _, name := range names {

			cost := synn.GasCost(name)
			if err := synn.RegisterGasMetadata(name, cost, entry.category, entry.description); err != nil {
				return fmt.Errorf("register gas metadata %s: %w", name, err)
			}
		}
	}

	return nil
=======
	register("consensus", "Core consensus lifecycle operations", "MineBlock")
	register("dao", "DAO creation and authority renewal", "CreateDAO", "UpdateMemberRole", "RenewAuthorityTerm")
	register("cross-chain", "Stage 24 cross-chain operations", "RegisterBridge", "BridgeDeposit", "BridgeClaim", "OpenConnection", "CloseConnection", "LockMint", "BurnRelease")
	register("node", "Stage 25 node and infrastructure operations", "SetMode", "Stake", "Unstake", "Optimize", "SecureCommand", "TrackLogistics", "ShareTactical", "ReportFork", "Metrics")
	register("templates", "Stage 29 contract templates", "DeployTokenFaucetTemplate", "DeployStorageMarketTemplate", "DeployDAOGovernanceTemplate", "DeployNFTMintingTemplate", "DeployAIModelMarketTemplate")
	register("marketplace", "Stage 34 marketplace settlement", "DeploySmartContract", "TradeContract")
	register("storage", "Stage 35 storage marketplace operations", "CreateListing", "ListListings", "GetListing", "OpenDeal", "CloseDeal", "ListDeals", "GetDeal", "Storage_Pin", "Storage_Retrieve", "IPFS_Add", "IPFS_Get", "IPFS_Unpin")
	register("nft", "Stage 36 NFT marketplace operations", "MintNFT", "ListNFT", "BuyNFT")
	register("dex", "Stage 39 liquidity view operations", "Liquidity_Pool", "Liquidity_Pools")
	register("wallet", "Wallet lifecycle operations", "NewWallet", "Sign", "VerifySignature")
	register("content", "Stage 59 content registry operations", "RegisterContentNode", "UploadContent", "RetrieveContent", "ListContentNodes")
	register("monetary", "Stage 40 monetary policy queries", "BlockReward", "CirculatingSupply", "RemainingSupply", "InitialPrice", "AlphaFactor", "MinimumStake")
	register("p2p", "Stage 67 Kademlia routing operations", "KademliaStore", "KademliaGet", "KademliaClosest", "KademliaDistance")
	register("orchestrator", "Stage 78 enterprise orchestrator operations", "EnterpriseBootstrap", "EnterpriseConsensusSync", "EnterpriseWalletSeal", "EnterpriseNodeAudit", "EnterpriseAuthorityElect")
	register("enterprise", "Stage 79 enterprise combined node operations", "EnterpriseSpecialAttach", "EnterpriseSpecialDetach", "EnterpriseSpecialBroadcast", "EnterpriseSpecialSnapshot", "EnterpriseSpecialLedger")
	logrus.Debug("gas table loaded")

	// Preload stage 3 modules so CLI commands can operate without extra setup.
	_ = core.NewAuthorityNodeRegistry()
	_ = core.NewBankInstitutionalNode("init", "init", core.NewLedger())

	// Preload stage 8 modules to expose contract and cross-chain managers via CLI.
	vm := core.NewSimpleVM()
	_ = vm.Start()
	_ = core.NewContractRegistry(vm)
	_ = core.NewBridgeRegistry()
	_ = core.NewBridgeTransferManager()
	_ = core.NewChainConnectionManager()
	_ = core.NewProtocolRegistry()
	_ = core.NewCrossChainTxManager(core.NewLedger())

	// Preload stage 11 modules for VM sandbox management.
	_ = core.NewSandboxManager()

	// Preload stage 9 modules so DAO-related CLI commands are ready for use.
	daoMgr := core.NewDAOManager()
	_ = core.NewProposalManager()
	_ = core.NewDAOStaking(daoMgr)
	_ = core.NewDAOTokenLedger(daoMgr)
	_ = core.NewConsensusNetworkManager()
	_ = core.NewCustodialNode("cli-custodian", "cli-custodian", core.NewLedger())

	// Preload stage 12 modules to expose wallet, warfare, watchtower and data distribution monitoring
	// functionality via the CLI.
	if _, err := core.NewWallet(); err != nil {
		logrus.Debugf("wallet init error: %v", err)
	}
	_ = core.NewWarfareNode(core.NewNode("cli-war", "cli-war", core.NewLedger()))
	_ = core.NewWatchtowerNode("cli-watchtower", nil)
	// Stage 59 modules for content registry and secrets management
	_ = core.NewContentNetworkNode("cli-content", "cli")
	_ = security.NewSecretsManager()

	// Preload stage 13 modules for secure channels and compliance checks.
	_ = core.NewZeroTrustEngine()
	_ = core.NewRegulatoryNode("cli-regnode", core.NewRegulatoryManager())

	// Preload stage 20 token extensions for CLI and opcode availability.
	_ = tokens.NewSYN223Token("cli", "S223", "cli", 0)
	_ = tokens.NewSYN2700Token()
	_ = tokens.NewSYN3200Token(1)
	_ = tokens.NewSYN3600Token()
	_ = tokens.NewSYN3800Token(0)
	_ = tokens.NewSYN3900Token()
	_ = tokens.NewSYN500Token()
	_ = tokens.NewSYN5000Token()
	// Preload stage 36 NFT marketplace
	_ = core.NewNFTMarketplace()

	logrus.Infof("starting Synnergy in %s mode on %s:%d", cfg.Environment, cfg.Server.Host, cfg.Server.Port)

	if err := cli.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
