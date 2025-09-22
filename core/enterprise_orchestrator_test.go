package core

import (
	"context"
	"testing"
)

func TestEnterpriseOrchestratorDiagnosticsUnit(t *testing.T) {
	orch, err := NewEnterpriseOrchestrator(context.Background())
	if err != nil {
		t.Fatalf("unexpected error creating orchestrator: %v", err)
	}
	diag := orch.Diagnostics(context.Background())
	if !diag.VMRunning {
		t.Fatalf("expected VM to be running")
	}
	if diag.VMMode == "" {
		t.Fatalf("vm mode not reported")
	}
	if diag.WalletAddress == "" {
		t.Fatalf("wallet address missing")
	}
	if len(diag.GasCoverage) == 0 {
		t.Fatalf("gas coverage not populated")
	}
	for _, opcode := range []string{"EnterpriseBootstrap", "EnterpriseConsensusSync", "EnterpriseWalletSeal", "EnterpriseNodeAudit", "EnterpriseAuthorityElect"} {
		if diag.GasCoverage[opcode] == 0 {
			t.Fatalf("opcode %s missing from diagnostics", opcode)
		}
	}
	if diag.MissingOpcodes == nil {
		t.Fatalf("expected diagnostics to include missing opcode slice")
	}
}

func TestEnterpriseOrchestratorConsensusSituational(t *testing.T) {
	orch, err := NewEnterpriseOrchestrator(context.Background())
	if err != nil {
		t.Fatalf("unexpected error creating orchestrator: %v", err)
	}
	id, err := orch.RegisterConsensusNetwork(context.Background(), "ProofOfStake", "Synnergy-PBFT")
	if err != nil {
		t.Fatalf("failed to register consensus network: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected non-zero network id")
	}
	diag := orch.Diagnostics(context.Background())
	if diag.ConsensusNetworks == 0 {
		t.Fatalf("consensus network count not reflected in diagnostics")
	}
}

func TestEnterpriseOrchestratorGasStress(t *testing.T) {
	orch, err := NewEnterpriseOrchestrator(context.Background())
	if err != nil {
		t.Fatalf("unexpected error creating orchestrator: %v", err)
	}
	for i := 0; i < 5; i++ {
		diag, err := orch.SyncGasSchedule(context.Background(), map[string]uint64{"EnterpriseWalletSeal": 60})
		if err != nil {
			t.Fatalf("sync gas failed: %v", err)
		}
		if diag.GasCoverage["EnterpriseWalletSeal"] != 60 {
			t.Fatalf("gas schedule not applied, got %d", diag.GasCoverage["EnterpriseWalletSeal"])
		}
	}
}

func TestEnterpriseOrchestratorFunctionalAuthority(t *testing.T) {
	orch, err := NewEnterpriseOrchestrator(context.Background())
	if err != nil {
		t.Fatalf("unexpected error creating orchestrator: %v", err)
	}
	node, err := orch.RegisterAuthorityNode(context.Background(), "stage78-operator", "governor")
	if err != nil {
		t.Fatalf("register authority node failed: %v", err)
	}
	if node.Address != "stage78-operator" {
		t.Fatalf("unexpected authority node address: %s", node.Address)
	}
	diag := orch.Diagnostics(context.Background())
	if diag.AuthorityNodes < 2 {
		t.Fatalf("expected authority nodes to include orchestrator wallet and new node")
	}
}

func TestEnterpriseOrchestratorRealWorldFlow(t *testing.T) {
	orch, err := NewEnterpriseOrchestrator(context.Background())
	if err != nil {
		t.Fatalf("unexpected error creating orchestrator: %v", err)
	}
	if _, err := orch.RegisterConsensusNetwork(context.Background(), "SYN-PBFT", "External-PoS"); err != nil {
		t.Fatalf("real world network registration failed: %v", err)
	}
	if _, err := orch.RegisterAuthorityNode(context.Background(), "stage78-delegate", "compliance"); err != nil {
		t.Fatalf("real world authority registration failed: %v", err)
	}
	diag, err := orch.SyncGasSchedule(context.Background(), map[string]uint64{"EnterpriseBootstrap": 120})
	if err != nil {
		t.Fatalf("sync gas during real world flow failed: %v", err)
	}
	if diag.ConsensusNetworks == 0 {
		t.Fatalf("expected consensus network count in diagnostics")
	}
	if diag.AuthorityNodes < 2 {
		t.Fatalf("expected authority node count to include additional member")
	}
	if diag.WalletAddress == "" {
		t.Fatalf("wallet address missing from diagnostics")
	}
}
