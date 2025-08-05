package synnergy

import "testing"

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
	id := m.OpenConnection("chainA", "chainB")
	if len(m.ListConnections()) != 1 {
		t.Fatalf("expected 1 connection")
	}
	if err := m.CloseConnection(id); err != nil {
		t.Fatalf("close failed: %v", err)
	}
	if err := m.CloseConnection(id); err == nil {
		t.Fatalf("expected error on double close")
	}
}

func TestXContractRegistry(t *testing.T) {
	r := NewXContractRegistry()
	r.RegisterMapping("local1", "chainB", "remote1")
	if len(r.ListMappings()) != 1 {
		t.Fatalf("expected 1 mapping")
	}
	m, ok := r.GetMapping("local1")
	if !ok || m.RemoteAddress != "remote1" {
		t.Fatalf("unexpected mapping: %+v", m)
	}
	r.RemoveMapping("local1")
	if _, ok := r.GetMapping("local1"); ok {
		t.Fatalf("mapping should be removed")
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
