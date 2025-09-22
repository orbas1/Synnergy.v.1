package main

import (
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	synn "synnergy"
	"synnergy/cli"
	"synnergy/core"
	"synnergy/internal/config"
	security "synnergy/internal/security"
	tokens "synnergy/internal/tokens"
)

func main() {
	// Load variables from .env if present to mirror the setup guides.
	gotenv.Load()

	// Initialize a no-op tracer provider so modules can emit spans safely.
	otel.SetTracerProvider(trace.NewNoopTracerProvider())

	cfgPath := os.Getenv("SYN_CONFIG")
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		logrus.Fatal(err)
	}

	// Configure logging based on loaded configuration.
	lvl, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		logrus.Fatalf("invalid log level: %v", err)
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

	// Warm up caches for shared resources and ensure gas costs are registered.
	synn.LoadGasTable()
	stage78Gas := map[string]uint64{
		"EnterpriseBootstrap":      120,
		"EnterpriseConsensusSync":  95,
		"EnterpriseWalletSeal":     60,
		"EnterpriseNodeAudit":      75,
		"EnterpriseAuthorityElect": 80,
	}
	if inserted, err := synn.EnsureGasSchedule(stage78Gas); err != nil {
		logrus.Fatalf("stage 78 gas sync failed: %v", err)
	} else if len(inserted) > 0 {
		logrus.Infof("registered %d stage 78 opcodes", len(inserted))
	}
	stage81Gas := map[string]uint64{
		"ModuleCatalogueList":    3,
		"ModuleCatalogueInspect": 5,
		"WalletNew":              20,
		"VMCreate":               15,
		"VMStart":                10,
		"VMStop":                 6,
		"VMStatus":               4,
		"VMExec":                 45,
		"NodeInfo":               4,
		"NodeStake":              25,
		"NodeSlash":              30,
		"NodeRehab":              12,
		"NodeAddTx":              8,
		"NodeMempool":            3,
		"NodeMine":               55,
	}
	if inserted, err := synn.EnsureGasSchedule(stage81Gas); err != nil {
		logrus.Fatalf("stage 81 gas sync failed: %v", err)
	} else if len(inserted) > 0 {
		logrus.Infof("registered %d stage 81 opcodes", len(inserted))
	}
	register := func(category, description string, names ...string) {
		for _, name := range names {
			cost := synn.GasCost(name)
			if err := synn.RegisterGasMetadata(name, cost, category, description); err != nil {
				logrus.Fatalf("register gas metadata %s: %v", name, err)
			}
		}
	}
	register("consensus", "Core consensus lifecycle operations", "MineBlock")
	register("dao", "DAO creation and authority renewal", "CreateDAO", "UpdateMemberRole", "RenewAuthorityTerm")
	register("cross-chain", "Stage 24 cross-chain operations", "RegisterBridge", "BridgeDeposit", "BridgeClaim", "OpenConnection", "CloseConnection", "LockMint", "BurnRelease")
	register("node", "Stage 25 node and infrastructure operations", "SetMode", "Stake", "Unstake", "Optimize", "SecureCommand", "TrackLogistics", "ShareTactical", "ReportFork", "Metrics", "NodeInfo", "NodeStake", "NodeSlash", "NodeRehab", "NodeAddTx", "NodeMempool", "NodeMine")
	register("templates", "Stage 29 contract templates", "DeployTokenFaucetTemplate", "DeployStorageMarketTemplate", "DeployDAOGovernanceTemplate", "DeployNFTMintingTemplate", "DeployAIModelMarketTemplate")
	register("marketplace", "Stage 34 marketplace settlement", "DeploySmartContract", "TradeContract")
	register("storage", "Stage 35 storage marketplace operations", "CreateListing", "ListListings", "GetListing", "OpenDeal", "CloseDeal", "ListDeals", "GetDeal", "Storage_Pin", "Storage_Retrieve", "IPFS_Add", "IPFS_Get", "IPFS_Unpin")
	register("nft", "Stage 36 NFT marketplace operations", "MintNFT", "ListNFT", "BuyNFT")
	register("dex", "Stage 39 liquidity view operations", "Liquidity_Pool", "Liquidity_Pools")
	register("wallet", "Wallet lifecycle operations", "NewWallet", "Sign", "VerifySignature", "WalletNew")
	register("content", "Stage 59 content registry operations", "RegisterContentNode", "UploadContent", "RetrieveContent", "ListContentNodes")
	register("monetary", "Stage 40 monetary policy queries", "BlockReward", "CirculatingSupply", "RemainingSupply", "InitialPrice", "AlphaFactor", "MinimumStake")
	register("p2p", "Stage 67 Kademlia routing operations", "KademliaStore", "KademliaGet", "KademliaClosest", "KademliaDistance")
	register("orchestrator", "Stage 78 enterprise orchestrator operations", "EnterpriseBootstrap", "EnterpriseConsensusSync", "EnterpriseWalletSeal", "EnterpriseNodeAudit", "EnterpriseAuthorityElect")
	register("vm", "Stage 81 VM lifecycle operations", "VMCreate", "VMStart", "VMStop", "VMStatus", "VMExec")
	register("cli", "Stage 81 CLI module catalogue operations", "ModuleCatalogueList", "ModuleCatalogueInspect")
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
