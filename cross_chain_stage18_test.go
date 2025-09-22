package synnergy

import (
	"context"
	"errors"
	"testing"
)

func TestCrossChainManager(t *testing.T) {
	m := NewCrossChainManager()
	id := m.RegisterBridge("chainA", "chainB", "relayer1")
	if id == "" {
		t.Fatalf("expected bridge ID")
	}
	if len(m.ListBridges()) != 1 {
		t.Fatalf("expected 1 bridge")
	}
	b, ok := m.GetBridge(id)
	if !ok || b.SourceChain != "chainA" || b.TargetChain != "chainB" {
		t.Fatalf("unexpected bridge data: %+v", b)
	}
	if m.IsRelayerAuthorized("relayer1") {
		t.Fatalf("relayer should not be globally authorized yet")
	}
	m.AuthorizeRelayer("relayer1")
	if !m.IsRelayerAuthorized("relayer1") {
		t.Fatalf("relayer should be authorized")
	}
	m.RevokeRelayer("relayer1")
	if m.IsRelayerAuthorized("relayer1") {
		t.Fatalf("relayer should be revoked")
	}
}

func TestProtocolRegistry(t *testing.T) {
	r := NewProtocolRegistry()
	id := r.RegisterProtocol("IBC")
	if id == "" {
		t.Fatalf("expected protocol ID")
	}
	if len(r.ListProtocols()) != 1 {
		t.Fatalf("expected 1 protocol")
	}
	p, ok := r.GetProtocol(id)
	if !ok || p.Name != "IBC" {
		t.Fatalf("unexpected protocol: %+v", p)
	}
}

func TestBridgeTransferManager(t *testing.T) {
	m := NewBridgeTransferManager()
	id := m.Deposit("bridge1", "alice", "bob", 10, "tokenX")
	if id == "" {
		t.Fatalf("expected transfer ID")
	}
	if len(m.ListTransfers()) != 1 {
		t.Fatalf("expected 1 transfer")
	}
	if err := m.Claim(id, []byte("proof")); err != nil {
		t.Fatalf("claim failed: %v", err)
	}
	if err := m.Claim(id, []byte("proof")); err == nil {
		t.Fatalf("expected error on double claim")
	}
}

func TestConnectionManager(t *testing.T) {
	m := NewConnectionManager()
	ctx := context.Background()
	conn, err := m.OpenConnection(ctx, ConnectionSpec{
		LocalChain:       "chainA",
		RemoteChain:      "chainB",
		LocalEndpoint:    "chainA-endpoint",
		RemoteEndpoint:   "chainB-endpoint",
		Signer:           "authority-chainA",
		HandshakeProof:   []byte("proof"),
		HandshakePayload: []byte("payload"),
	})
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	conns := m.ListConnections(ConnectionFilter{})
	if len(conns) != 1 {
		t.Fatalf("expected 1 connection")
	}
	closed, err := m.CloseConnection(ctx, conn.ID, "test shutdown")
	if err != nil {
		t.Fatalf("close failed: %v", err)
	}
	if closed.Status != ConnectionStatusClosed {
		t.Fatalf("expected closed status, got %s", closed.Status)
	}
	if _, err := m.CloseConnection(ctx, conn.ID, "double close"); !errors.Is(err, ErrConnectionClosed) {
		t.Fatalf("expected ErrConnectionClosed, got %v", err)
	}
}

func TestContractLinkManager(t *testing.T) {
	manager := NewContractLinkManager()
	spec := ContractLinkSpec{
		LocalChain:    "chainA",
		LocalAddress:  "0xabc",
		RemoteChain:   "chainB",
		RemoteAddress: "0xdef",
		ConnectionID:  "conn-1",
		Capabilities:  []string{"invoke"},
		GasProfile:    "balanced",
		Metadata: map[string]string{
			"department": "operations",
		},
		AccessPolicy: AccessPolicy{
			AllowedApprovers:  []string{"ops"},
			RequiredApprovals: 1,
		},
	}
	link, err := manager.Register(spec)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if link.Status != ContractLinkStatusPending {
		t.Fatalf("expected pending status, got %s", link.Status)
	}

	active, err := manager.RecordApproval(link.ID, "ops")
	if err != nil {
		t.Fatalf("approval failed: %v", err)
	}
	if active.Status != ContractLinkStatusActive {
		t.Fatalf("expected active status after approval, got %s", active.Status)
	}

	updated, err := manager.Update(link.ID, ContractLinkUpdate{
		GasProfile: "priority",
		Metadata: map[string]string{
			"department": "operations",
			"region":     "eu",
		},
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Spec.Metadata["region"] != "eu" {
		t.Fatalf("expected metadata update to persist, got %+v", updated.Spec.Metadata)
	}

	failed, err := manager.ReportFailure(link.ID, "HEARTBEAT_TIMEOUT", "remote contract unresponsive")
	if err != nil {
		t.Fatalf("report failure failed: %v", err)
	}
	if failed.Status != ContractLinkStatusFailed {
		t.Fatalf("expected failed status, got %s", failed.Status)
	}

	resumed, err := manager.Resume(link.ID)
	if err != nil {
		t.Fatalf("resume failed: %v", err)
	}
	if resumed.Status != ContractLinkStatusActive {
		t.Fatalf("expected active status after resume, got %s", resumed.Status)
	}

	suspended, err := manager.Suspend(link.ID, "governance review")
	if err != nil {
		t.Fatalf("suspend failed: %v", err)
	}
	if suspended.Status != ContractLinkStatusSuspended {
		t.Fatalf("expected suspended status, got %s", suspended.Status)
	}

	resumed, err = manager.Resume(link.ID)
	if err != nil {
		t.Fatalf("resume after suspension failed: %v", err)
	}
	if resumed.Status != ContractLinkStatusActive {
		t.Fatalf("expected active status, got %s", resumed.Status)
	}

	retired, err := manager.Retire(link.ID, "sunset")
	if err != nil {
		t.Fatalf("retire failed: %v", err)
	}
	if retired.Status != ContractLinkStatusRetired {
		t.Fatalf("expected retired status, got %s", retired.Status)
	}

	if matches := manager.List(ContractLinkFilter{ConnectionID: "conn-1", Statuses: []ContractLinkStatus{ContractLinkStatusRetired}}); len(matches) != 1 {
		t.Fatalf("expected retired mapping to be discoverable, got %d", len(matches))
	}

	spec.RemoteAddress = "0x999"
	spec.Metadata["department"] = "innovation"
	replacement, err := manager.Register(spec)
	if err != nil {
		t.Fatalf("expected replacement mapping to register, got %v", err)
	}
	if replacement.ID == link.ID {
		t.Fatalf("replacement should allocate a new identifier")
	}
}

func TestTransactionManager(t *testing.T) {
	m := NewTransactionManager()
	id1 := m.LockAndMint("bridge1", "asset", 5, "proof")
	if id1 == "" {
		t.Fatalf("expected tx id")
	}
	id2 := m.BurnAndRelease("bridge1", "bob", "asset", 5)
	if len(m.ListTransactions()) != 2 {
		t.Fatalf("expected 2 transactions")
	}
	if _, ok := m.GetTransaction(id2); !ok {
		t.Fatalf("transaction not found")
	}
}
