package main

import (
	"os"

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
	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Fatalf("invalid log level: %v", err)
	}
	logrus.SetLevel(lvl)

	// Warm up caches for shared resources and ensure gas costs are registered.
	synn.LoadGasTable()
	mustRegister := func(name string) {
		if err := synn.RegisterGasCost(name, synn.GasCost(name)); err != nil {
			logrus.Fatalf("register gas cost %s: %v", name, err)
		}
	}
	mustRegister("MineBlock")
	mustRegister("CreateDAO")
	mustRegister("UpdateMemberRole")
	mustRegister("RenewAuthorityTerm")
	// Stage 24 cross-chain operations
	mustRegister("RegisterBridge")
	mustRegister("BridgeDeposit")
	mustRegister("BridgeClaim")
	mustRegister("OpenConnection")
	mustRegister("CloseConnection")
	mustRegister("LockMint")
	mustRegister("BurnRelease")
	// Stage 25 node and operations costs
	mustRegister("SetMode")
	mustRegister("Stake")
	mustRegister("Unstake")
	mustRegister("Optimize")
	mustRegister("SecureCommand")
	mustRegister("TrackLogistics")
	mustRegister("ShareTactical")
	mustRegister("ReportFork")
	mustRegister("Metrics")
	// Stage 29 smart contract templates
	mustRegister("DeployTokenFaucetTemplate")
	mustRegister("DeployStorageMarketTemplate")
	mustRegister("DeployDAOGovernanceTemplate")
	mustRegister("DeployNFTMintingTemplate")
	mustRegister("DeployAIModelMarketTemplate")
	// Stage 34 smart-contract marketplace operations
	mustRegister("DeploySmartContract")
	mustRegister("TradeContract")
	// Stage 35 storage marketplace operations
	mustRegister("CreateListing")
	mustRegister("ListListings")
	mustRegister("GetListing")
	mustRegister("OpenDeal")
	mustRegister("CloseDeal")
	mustRegister("ListDeals")
	mustRegister("GetDeal")
	mustRegister("Storage_Pin")
	mustRegister("Storage_Retrieve")
	mustRegister("IPFS_Add")
	mustRegister("IPFS_Get")
	mustRegister("IPFS_Unpin")
	// Stage 36 NFT marketplace operations
	mustRegister("MintNFT")
	mustRegister("ListNFT")
	mustRegister("BuyNFT")
	// Stage 39 liquidity view operations for DEX screener
	mustRegister("Liquidity_Pool")
	mustRegister("Liquidity_Pools")
	// Wallet operations used by GUI clients
	mustRegister("NewWallet")
	mustRegister("Sign")
	mustRegister("VerifySignature")
	// Stage 59 content node management
	mustRegister("RegisterContentNode")
	mustRegister("UploadContent")
	mustRegister("RetrieveContent")
	mustRegister("ListContentNodes")
	// Stage 40 monetary policy queries
	mustRegister("BlockReward")
	mustRegister("CirculatingSupply")
	mustRegister("RemainingSupply")
	mustRegister("InitialPrice")
	mustRegister("AlphaFactor")
	mustRegister("MinimumStake")
	// Stage 67 Kademlia operations
	mustRegister("KademliaStore")
	mustRegister("KademliaGet")
	mustRegister("KademliaClosest")
	mustRegister("KademliaDistance")
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
