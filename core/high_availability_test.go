package core

import (
	"context"
	"testing"
	"time"
)

// TestFailoverManagerKeepsPrimaryWhenHealthy ensures that the configured
// primary node remains active as long as it has sent a recent heartbeat.
func TestFailoverManagerKeepsPrimaryWhenHealthy(t *testing.T) {
	timeout := 50 * time.Millisecond
	fm := NewFailoverManager("primary", timeout)
	fm.RegisterBackup("backup")

	if active := fm.Active(); active != "primary" {
		t.Fatalf("expected primary to remain active, got %s", active)
	}

	if err := fm.RecordHeartbeat(HeartbeatProof{ID: "backup"}); err != nil {
		t.Fatalf("heartbeat: %v", err)
	}
	if active := fm.Active(); active != "primary" {
		t.Fatalf("expected primary after backup heartbeat, got %s", active)
	}
}

// TestFailoverManagerFailoverToLatestBackup verifies that when the primary
// fails to heartbeat within the timeout, the most recently heartbeating backup
// is promoted.
func TestFailoverManagerFailoverToLatestBackup(t *testing.T) {
	timeout := 5 * time.Millisecond
	fm := NewFailoverManager("p1", timeout)
	fm.RegisterBackup("b1")
	fm.RegisterBackup("b2")

	fm.mu.Lock()
	fm.nodes["p1"].LastHeartbeat = time.Now().Add(-2 * timeout)
	fm.nodes["b1"].LastHeartbeat = time.Now().Add(-50 * time.Millisecond)
	fm.nodes["b2"].LastHeartbeat = time.Now()
	fm.mu.Unlock()

	if active := fm.Active(); active != "b2" {
		t.Fatalf("expected b2 to become primary, got %s", active)
	}
	if active := fm.Active(); active != "b2" {
		t.Fatalf("expected b2 to remain primary after promotion, got %s", active)
	}
}

// TestFailoverManagerHeartbeatAndRegister tests that registering a backup and
// issuing heartbeats correctly updates internal state.
func TestFailoverManagerHeartbeatAndRegister(t *testing.T) {
	timeout := 20 * time.Millisecond
	fm := NewFailoverManager("p1", timeout)

	fm.mu.RLock()
	primaryHB := fm.nodes["p1"].LastHeartbeat
	fm.mu.RUnlock()

	time.Sleep(1 * time.Millisecond)
	if err := fm.RecordHeartbeat(HeartbeatProof{ID: "p1"}); err != nil {
		t.Fatalf("heartbeat: %v", err)
	}

	fm.mu.RLock()
	updatedHB := fm.nodes["p1"].LastHeartbeat
	fm.mu.RUnlock()
	if !updatedHB.After(primaryHB) {
		t.Fatalf("heartbeat did not update timestamp")
	}

	fm.RegisterBackup("b1")
	if err := fm.RecordHeartbeat(HeartbeatProof{ID: "b1"}); err != nil {
		t.Fatalf("heartbeat backup: %v", err)
	}
	fm.mu.RLock()
	_, ok := fm.nodes["b1"]
	fm.mu.RUnlock()
	if !ok {
		t.Fatalf("backup node not registered")
	}
}

// TestFailoverManagerRemoveNode ensures nodes can be removed and primaries are
// re-elected if necessary.
func TestFailoverManagerRemoveNode(t *testing.T) {
	timeout := 10 * time.Millisecond
	fm := NewFailoverManager("p1", timeout)
	fm.RegisterBackup("b1")
	fm.RegisterBackup("b2")
	_ = fm.RecordHeartbeat(HeartbeatProof{ID: "b1"})
	_ = fm.RecordHeartbeat(HeartbeatProof{ID: "b2"})

	fm.RemoveNode("p1")
	if active := fm.Active(); active == "p1" || active == "" {
		t.Fatalf("expected backup to be promoted after removal, got %s", active)
	}

	fm.RemoveNode("b1")
	fm.RemoveNode("b2")
	if active := fm.Active(); active != "" {
		t.Fatalf("expected no active node after removing all, got %s", active)
	}
}

// TestFailoverManagerReportIntegrations validates the resilience report surfaces
// VM, consensus, wallet, ledger and audit data.
func TestFailoverManagerReportIntegrations(t *testing.T) {
	wallet, err := NewWallet()
	if err != nil {
		t.Fatalf("wallet init: %v", err)
	}
	fm := NewFailoverManager(wallet.Address, 25*time.Millisecond,
		WithFailoverVirtualMachine(NewSimpleVM(VMLight)),
		WithFailoverConsensus(NewConsensusNetworkManager()),
		WithFailoverWallet(wallet),
		WithFailoverRegistry(NewAuthorityNodeRegistry()),
		WithFailoverLedger(NewLedger()),
	)
	fm.RegisterBackup("backup")
	if err := fm.RecordHeartbeat(HeartbeatProof{ID: "backup", Latency: 10 * time.Millisecond}); err != nil {
		t.Fatalf("heartbeat: %v", err)
	}

	diag := fm.Report(context.Background())
	if diag.ActiveNode != wallet.Address {
		t.Fatalf("expected orchestrator wallet to remain active, got %s", diag.ActiveNode)
	}
	if diag.WalletAddress == "" || !diag.WalletReady {
		t.Fatalf("wallet readiness not surfaced: %+v", diag)
	}
	if diag.Signature == "" {
		t.Fatalf("report signature missing")
	}
	if len(diag.AuditTrail) == 0 {
		t.Fatalf("expected audit trail entries")
	}
}

// TestFailoverManagerRecordHeartbeatVerification ensures digital signatures are enforced.
func TestFailoverManagerRecordHeartbeatVerification(t *testing.T) {
	wallet, err := NewWallet()
	if err != nil {
		t.Fatalf("wallet init: %v", err)
	}
	fm := NewFailoverManager("primary", 10*time.Millisecond)
	fm.RegisterNode(FailoverNode{ID: "signed", Role: "validator", Region: "us-east", PublicKey: wallet.PublicKeyBytes()})

	payload := []byte("signed-heartbeat")
	sig, err := wallet.SignMessage(payload)
	if err != nil {
		t.Fatalf("sign message: %v", err)
	}
	if err := fm.RecordHeartbeat(HeartbeatProof{ID: "signed", Payload: payload, Signature: sig}); err != nil {
		t.Fatalf("record heartbeat: %v", err)
	}

	fm.mu.RLock()
	verified := fm.nodes["signed"].SignatureVerified
	fm.mu.RUnlock()
	if !verified {
		t.Fatalf("expected heartbeat signature verification to succeed")
	}
}
