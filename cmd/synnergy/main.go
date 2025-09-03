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

	// Preload stage 9 modules so DAO-related CLI commands are ready for use.
	_ = core.NewDAOManager()
	_ = core.NewProposalManager()
	_ = core.NewDAOStaking()
	_ = core.NewDAOTokenLedger()
	_ = core.NewConsensusNetworkManager()
	_ = core.NewCustodialNode("cli-custodian", "cli-custodian", core.NewLedger())

	logrus.Infof("starting Synnergy in %s mode on %s:%d", cfg.Environment, cfg.Server.Host, cfg.Server.Port)

	if err := cli.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
