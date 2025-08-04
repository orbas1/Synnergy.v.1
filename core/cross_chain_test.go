package core

import (
	"math/big"
	"testing"
)

func TestCrossChainModules(t *testing.T) {
	// Bridge registry
	br := NewBridgeRegistry()
	bridge, err := br.RegisterBridge("chainA", "chainB", "relayer1")
	if err != nil {
		t.Fatalf("register bridge: %v", err)
	}
	if _, ok := br.GetBridge(bridge.ID); !ok {
		t.Fatalf("bridge not found")
	}
	if err := br.AuthorizeRelayer(bridge.ID, "relayer2"); err != nil {
		t.Fatalf("authorize: %v", err)
	}
	if err := br.RevokeRelayer(bridge.ID, "relayer2"); err != nil {
		t.Fatalf("revoke: %v", err)
	}

	// Protocol registry
	pr := NewProtocolRegistry()
	proto, err := pr.RegisterProtocol("IBC")
	if err != nil {
		t.Fatalf("register protocol: %v", err)
	}
	if _, ok := pr.GetProtocol(proto.ID); !ok {
		t.Fatalf("protocol missing")
	}

	// Bridge transfer
	bt := NewBridgeTransferManager()
	tr, err := bt.Deposit(bridge.ID, "alice", "bob", big.NewInt(100), "token1")
	if err != nil {
		t.Fatalf("deposit: %v", err)
	}
	if err := bt.Claim(tr.ID, []byte("proof")); err != nil {
		t.Fatalf("claim: %v", err)
	}
	if !tr.Claimed {
		t.Fatalf("expected claimed transfer")
	}

	// Connection manager
	cm := NewConnectionManager()
	conn, err := cm.OpenConnection("chainA", "chainC")
	if err != nil {
		t.Fatalf("open connection: %v", err)
	}
	if err := cm.CloseConnection(conn.ID); err != nil {
		t.Fatalf("close connection: %v", err)
	}

	// Contract registry
	cr := NewContractRegistry()
	if _, err := cr.RegisterContract("0xabc", "chainB", "0xdef"); err != nil {
		t.Fatalf("register contract: %v", err)
	}
	if _, ok := cr.GetMapping("0xabc"); !ok {
		t.Fatalf("mapping missing")
	}
	if err := cr.RemoveMapping("0xabc"); err != nil {
		t.Fatalf("remove mapping: %v", err)
	}

	// Cross-chain transactions
	txm := NewCrossChainTxManager()
	tx1, err := txm.LockAndMint(bridge.ID, "asset1", big.NewInt(50), "proof")
	if err != nil {
		t.Fatalf("lockmint: %v", err)
	}
	if _, ok := txm.GetTx(tx1.ID); !ok {
		t.Fatalf("tx not found")
	}
	if _, err := txm.BurnAndRelease(bridge.ID, "bob", "asset1", big.NewInt(20)); err != nil {
		t.Fatalf("burnrelease: %v", err)
	}
	if len(txm.ListTxs()) != 2 {
		t.Fatalf("expected 2 transactions")
	}

	// Cross-consensus networks
	ccs := NewCCSNetworkRegistry()
	net, err := ccs.RegisterNetwork("PoS", "PoW")
	if err != nil {
		t.Fatalf("register network: %v", err)
	}
	if _, ok := ccs.GetNetwork(net.ID); !ok {
		t.Fatalf("network not found")
	}
}
