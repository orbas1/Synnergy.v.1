package cli

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"synnergy/core"
)

// TestDAOStakingStakeJSON verifies staking with JSON output and signature verification.
func TestDAOStakingStakeJSON(t *testing.T) {
	daoMgr = core.NewDAOManager()
	ledger = core.NewLedger()
	ledger.Mint("addr1", 100)
	daoStaking = core.NewDAOStaking(daoMgr, ledger)
	daoMgr.AuthorizeRelayer("addr1")
	dao, err := daoMgr.Create("dao", "addr1")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	addr := "addr1"
	amt := "10"
	msg := []byte(dao.ID + addr + amt)
	h := sha256.Sum256(msg)
	r, s, err := ecdsa.Sign(rand.Reader, priv, h[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	sig := append(r.Bytes(), s.Bytes()...)
	sigHex := hex.EncodeToString(sig)
	pub := append([]byte{4}, priv.PublicKey.X.Bytes()...)
	pub = append(pub, priv.PublicKey.Y.Bytes()...)
	pubHex := hex.EncodeToString(pub)
	msgHex := hex.EncodeToString(msg)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"dao-stake", "stake", dao.ID, addr, amt, "--pub", pubHex, "--msg", msgHex, "--sig", sigHex, "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	var resp map[string]string
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["status"] != "staked" {
		t.Fatalf("unexpected status: %v", resp)
	}
	if daoStaking.Balance(dao.ID, addr) != 10 {
		t.Fatalf("stake not recorded")
	}
}
