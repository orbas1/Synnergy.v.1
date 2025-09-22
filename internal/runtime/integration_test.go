package runtime

import (
	"context"
	"testing"
	"time"

	"synnergy/core"
)

func TestRuntimeIntegrationLifecycle(t *testing.T) {
	cfg := IntegrationConfig{
		NodeID:              "stage77-node",
		NodeAddress:         "127.0.0.1:9000",
		AuthorityAddress:    "authority-stage77",
		AuthorityDepartment: "compliance",
		CriticalOpcodes: []string{
			"MineBlock",
			"RegisterContentNode",
		},
		GasLimit:        64,
		MonitorInterval: 20 * time.Millisecond,
	}
	ri, err := NewRuntimeIntegration(cfg)
	if err != nil {
		t.Fatalf("new integration: %v", err)
	}
	if err := ri.EnsureGasSchedule(); err != nil {
		t.Fatalf("ensure gas schedule: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := ri.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer ri.Stop()

	ri.Ledger().Credit(ri.WalletAddress(), 100)
	tx := core.NewTransaction(ri.WalletAddress(), "recipient", 10, 1, 1)
	if err := ri.SubmitTransaction(context.Background(), tx); err != nil {
		t.Fatalf("submit transaction: %v", err)
	}

	// Wait for a monitoring tick to occur.
	select {
	case <-ri.HealthChannel():
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("health channel did not emit snapshot")
	}

	health := ri.HealthReport()
	if !health.Started {
		t.Fatalf("expected integration marked as started")
	}
	if !health.VMRunning {
		t.Fatalf("expected VM running")
	}
	if health.WalletAddress == "" {
		t.Fatalf("wallet address missing")
	}
}

func TestRuntimeIntegrationEnsureGasScheduleError(t *testing.T) {
	cfg := IntegrationConfig{
		NodeID:           "stage77-node-2",
		NodeAddress:      "127.0.0.1:9001",
		AuthorityAddress: "authority-stage77b",
		CriticalOpcodes:  []string{"NonExistentOpcode"},
	}
	ri, err := NewRuntimeIntegration(cfg)
	if err != nil {
		t.Fatalf("new integration: %v", err)
	}
	if err := ri.EnsureGasSchedule(); err == nil {
		t.Fatalf("expected error for missing opcode")
	}
}

func TestRuntimeIntegrationExecuteProgram(t *testing.T) {
	cfg := IntegrationConfig{
		NodeID:           "stage77-node-3",
		NodeAddress:      "127.0.0.1:9002",
		AuthorityAddress: "authority-stage77c",
		GasLimit:         10,
	}
	ri, err := NewRuntimeIntegration(cfg)
	if err != nil {
		t.Fatalf("new integration: %v", err)
	}
	if err := ri.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer ri.Stop()

	wasm := []byte{0x00, 0x00, 0x00}
	if _, _, err := ri.ExecuteProgram(context.Background(), wasm, "", []byte("payload")); err != nil {
		t.Fatalf("execute program: %v", err)
	}
}
