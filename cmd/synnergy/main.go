package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	synn "synnergy"
	"synnergy/cli"
	"synnergy/core"
	"synnergy/internal/config"
	runtime "synnergy/internal/runtime"
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
	criticalOps := []string{
		"MineBlock",
		"CreateDAO",
		"UpdateMemberRole",
		"RenewAuthorityTerm",
		"RegisterBridge",
		"BridgeDeposit",
		"BridgeClaim",
		"OpenConnection",
		"CloseConnection",
		"LockMint",
		"BurnRelease",
		"SetMode",
		"Stake",
		"Unstake",
		"Optimize",
		"SecureCommand",
		"TrackLogistics",
		"ShareTactical",
		"ReportFork",
		"Metrics",
		"DeployTokenFaucetTemplate",
		"DeployStorageMarketTemplate",
		"DeployDAOGovernanceTemplate",
		"DeployNFTMintingTemplate",
		"DeployAIModelMarketTemplate",
		"DeploySmartContract",
		"TradeContract",
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
		"MintNFT",
		"ListNFT",
		"BuyNFT",
		"Liquidity_Pool",
		"Liquidity_Pools",
		"NewWallet",
		"Sign",
		"VerifySignature",
		"RegisterContentNode",
		"UploadContent",
		"RetrieveContent",
		"ListContentNodes",
		"BlockReward",
		"CirculatingSupply",
		"RemainingSupply",
		"InitialPrice",
		"AlphaFactor",
		"MinimumStake",
		"KademliaStore",
		"KademliaGet",
		"KademliaClosest",
		"KademliaDistance",
	}

	integration, err := runtime.NewRuntimeIntegration(runtime.IntegrationConfig{
		NodeID:              "cli-runtime-node",
		NodeAddress:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		AuthorityAddress:    "cli-authority",
		AuthorityDepartment: "operations",
		CriticalOpcodes:     criticalOps,
		GasLimit:            10_000,
		MonitorInterval:     2 * time.Second,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	if err := integration.EnsureGasSchedule(); err != nil {
		logrus.Fatal(err)
	}
	if err := integration.Start(context.Background()); err != nil {
		logrus.Fatal(err)
	}
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

	go func() {
		for range integration.HealthChannel() {
			// Periodically log health snapshots so operators can
			// correlate CLI usage, VM workloads and consensus status.
			health := integration.HealthReport()
			logrus.WithFields(logrus.Fields{
				"node":       health.NodeID,
				"vm_running": health.VMRunning,
				"height":     health.LedgerHeight,
			}).Debug("runtime integration health")
		}
	}()

	if err := cli.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
