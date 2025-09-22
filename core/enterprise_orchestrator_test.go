package core

import (
	"context"
	"fmt"
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

func TestEnterpriseOrchestratorBootstrapUnit(t *testing.T) {
	orch, err := NewEnterpriseOrchestrator(context.Background())
	if err != nil {
		t.Fatalf("unexpected error creating orchestrator: %v", err)
	}
	res, err := orch.BootstrapNetwork(context.Background(), EnterpriseBootstrapConfig{
		NodeID:            "stage79-unit",
		EnableReplication: true,
		EnableRegulator:   true,
	})
	if err != nil {
		t.Fatalf("bootstrap unit config failed: %v", err)
	}
	if res.NodeID != "stage79-unit" {
		t.Fatalf("unexpected node id: %s", res.NodeID)
	}
	if res.ConsensusNetworkID == 0 {
		t.Fatalf("expected consensus network id")
	}
	if !res.ReplicationEnabled {
		t.Fatalf("replication flag not reported")
	}
	if res.BootstrapSignature == "" {
		t.Fatalf("bootstrap signature missing")
	}
	if !res.Diagnostics.ReplicationActive {
		t.Fatalf("diagnostics did not record replication activity")
	}
	if res.Diagnostics.BootstrapNodes == 0 {
		t.Fatalf("bootstrap node count not recorded")
	}
}

func TestEnterpriseOrchestratorBootstrapSituational(t *testing.T) {
	orch, err := NewEnterpriseOrchestrator(context.Background())
	if err != nil {
		t.Fatalf("unexpected error creating orchestrator: %v", err)
	}
	res, err := orch.BootstrapNetwork(context.Background(), EnterpriseBootstrapConfig{
		NodeID:            "stage79-situational",
		Address:           "127.0.0.1:9100",
		ConsensusProfile:  "External-PoS",
		GovernanceProfile: "SYN-DAO",
		Authorities: map[string]string{
			"stage79-authority": "compliance",
		},
	})
	if err != nil {
		t.Fatalf("bootstrap situational config failed: %v", err)
	}
	if res.Address != "127.0.0.1:9100" {
		t.Fatalf("unexpected address: %s", res.Address)
	}
	if len(res.AuthorityNodes) < 2 {
		t.Fatalf("expected authority nodes to include orchestrator and additional entry")
	}
	if res.Diagnostics.ConsensusNetworks == 0 {
		t.Fatalf("consensus networks not reflected in diagnostics")
	}
}

func TestEnterpriseOrchestratorBootstrapStress(t *testing.T) {
	orch, err := NewEnterpriseOrchestrator(context.Background())
	if err != nil {
		t.Fatalf("unexpected error creating orchestrator: %v", err)
	}
	for i := 0; i < 3; i++ {
		cfg := EnterpriseBootstrapConfig{NodeID: fmt.Sprintf("stage79-stress-%d", i)}
		if _, err := orch.BootstrapNetwork(context.Background(), cfg); err != nil {
			t.Fatalf("stress bootstrap %d failed: %v", i, err)
		}
	}
	diag := orch.Diagnostics(context.Background())
	if diag.BootstrapNodes < 3 {
		t.Fatalf("expected >=3 bootstrap nodes, got %d", diag.BootstrapNodes)
	}
}

func TestEnterpriseOrchestratorBootstrapFunctional(t *testing.T) {
	orch, err := NewEnterpriseOrchestrator(context.Background())
	if err != nil {
		t.Fatalf("unexpected error creating orchestrator: %v", err)
	}
	res, err := orch.BootstrapNetwork(context.Background(), EnterpriseBootstrapConfig{
		NodeID:          "stage79-functional",
		EnableRegulator: true,
		Authorities:     map[string]string{"stage79-governor": "governor"},
	})
	if err != nil {
		t.Fatalf("functional bootstrap failed: %v", err)
	}
	if res.Diagnostics.AuthorityNodes < 2 {
		t.Fatalf("expected orchestrator wallet and new authority to be registered")
	}
	if res.Diagnostics.WalletAddress == "" {
		t.Fatalf("wallet address missing in diagnostics")
	}
}

func TestEnterpriseOrchestratorBootstrapRealWorld(t *testing.T) {
	orch, err := NewEnterpriseOrchestrator(context.Background())
	if err != nil {
		t.Fatalf("unexpected error creating orchestrator: %v", err)
	}
	res, err := orch.BootstrapNetwork(context.Background(), EnterpriseBootstrapConfig{
		NodeID:            "stage79-realworld",
		EnableReplication: true,
		ConsensusProfile:  "SYN-PBFT",
		GovernanceProfile: "SYN-Gov",
		Authorities: map[string]string{
			"stage79-audit": "audit",
			"stage79-ops":   "operations",
		},
	})
	if err != nil {
		t.Fatalf("real world bootstrap failed: %v", err)
	}
	if res.BootstrapSignature == "" {
		t.Fatalf("expected bootstrap signature for audit trail")
	}
	if !res.Diagnostics.ReplicationActive {
		t.Fatalf("replication should be active for real world bootstrap")
	}
	if res.Diagnostics.ConsensusNetworks == 0 {
		t.Fatalf("expected consensus network count in diagnostics")
	}
}
