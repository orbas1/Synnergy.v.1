package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sort"
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

	enterpriseSpecialGas := map[string]uint64{
		"EnterpriseSpecialAttach":    110,
		"EnterpriseSpecialDetach":    55,
		"EnterpriseSpecialBroadcast": 145,
		"EnterpriseSpecialSnapshot":  40,
		"EnterpriseSpecialLedger":    30,
	}
	if inserted, err := synn.EnsureGasSchedule(enterpriseSpecialGas); err != nil {
		return fmt.Errorf("enterprise special gas sync failed: %w", err)
	} else if len(inserted) > 0 {
		logrus.Infof("registered %d enterprise special opcodes", len(inserted))
	}

	type category struct {
		name        string
		description string
		operations  []string
	}

	categories := []category{
		{name: "consensus", description: "Core consensus lifecycle operations", operations: []string{"MineBlock"}},
		{name: "dao", description: "DAO creation and authority renewal", operations: []string{"CreateDAO", "UpdateMemberRole", "RenewAuthorityTerm"}},
		{name: "cross-chain", description: "Stage 24 cross-chain operations", operations: []string{"RegisterBridge", "BridgeDeposit", "BridgeClaim", "OpenConnection", "CloseConnection", "LockMint", "BurnRelease"}},
		{name: "node", description: "Stage 25 node and infrastructure operations", operations: []string{"SetMode", "Stake", "Unstake", "Optimize", "SecureCommand", "TrackLogistics", "ShareTactical", "ReportFork", "Metrics"}},
		{name: "templates", description: "Stage 29 contract templates", operations: []string{"DeployTokenFaucetTemplate", "DeployStorageMarketTemplate", "DeployDAOGovernanceTemplate", "DeployNFTMintingTemplate", "DeployAIModelMarketTemplate"}},
		{name: "marketplace", description: "Stage 34 marketplace settlement", operations: []string{"DeploySmartContract", "TradeContract"}},
		{name: "storage", description: "Stage 35 storage marketplace operations", operations: []string{"CreateListing", "ListListings", "GetListing", "OpenDeal", "CloseDeal", "ListDeals", "GetDeal", "Storage_Pin", "Storage_Retrieve", "IPFS_Add", "IPFS_Get", "IPFS_Unpin"}},
		{name: "nft", description: "Stage 36 NFT marketplace operations", operations: []string{"MintNFT", "ListNFT", "BuyNFT"}},
		{name: "dex", description: "Stage 39 liquidity view operations", operations: []string{"Liquidity_Pool", "Liquidity_Pools"}},
		{name: "wallet", description: "Wallet lifecycle operations", operations: []string{"NewWallet", "Sign", "VerifySignature"}},
		{name: "content", description: "Stage 59 content registry operations", operations: []string{"RegisterContentNode", "UploadContent", "RetrieveContent", "ListContentNodes"}},
		{name: "monetary", description: "Stage 40 monetary policy queries", operations: []string{"BlockReward", "CirculatingSupply", "RemainingSupply", "InitialPrice", "AlphaFactor", "MinimumStake"}},
		{name: "p2p", description: "Stage 67 Kademlia routing operations", operations: []string{"KademliaStore", "KademliaGet", "KademliaClosest", "KademliaDistance"}},
		{name: "orchestrator", description: "Stage 78 enterprise orchestrator operations", operations: []string{"EnterpriseBootstrap", "EnterpriseConsensusSync", "EnterpriseWalletSeal", "EnterpriseNodeAudit", "EnterpriseAuthorityElect"}},
	}

	enterpriseOps := make([]string, 0, len(enterpriseSpecialGas))
	for op := range enterpriseSpecialGas {
		enterpriseOps = append(enterpriseOps, op)
	}
	sort.Strings(enterpriseOps)
	categories = append(categories, category{
		name:        "enterprise",
		description: "Stage 79 enterprise combined node operations",
		operations:  enterpriseOps,
	})

	for _, cat := range categories {
		for _, op := range cat.operations {
			cost := synn.GasCost(op)
			if err := synn.RegisterGasMetadata(op, cost, cat.name, cat.description); err != nil {
				return fmt.Errorf("register gas metadata %s: %w", op, err)
			}
		}
	}

	logrus.Debug("gas table loaded")
	return nil
}
