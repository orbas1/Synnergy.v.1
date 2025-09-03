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

	// Warm up caches for shared resources.
	synn.LoadGasTable()
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
	_ = core.NewDAOManager()
	_ = core.NewProposalManager()
	_ = core.NewDAOStaking()
	_ = core.NewDAOTokenLedger()
	_ = core.NewConsensusNetworkManager()
	_ = core.NewCustodialNode("cli-custodian", "cli-custodian", core.NewLedger())

	// Preload stage 12 modules to expose wallet, warfare and watchtower
	// functionality via the CLI.
	if _, err := core.NewWallet(); err != nil {
		logrus.Debugf("wallet init error: %v", err)
	}
	_ = core.NewWarfareNode(core.NewNode("cli-war", "cli-war", core.NewLedger()))
	_ = core.NewWatchtowerNode("cli-watchtower", nil)

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

	logrus.Infof("starting Synnergy in %s mode on %s:%d", cfg.Environment, cfg.Server.Host, cfg.Server.Port)

	if err := cli.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
