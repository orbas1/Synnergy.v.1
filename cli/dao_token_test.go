package cli

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"synnergy/core"
)

// TestDAOTokenMintJSON verifies mint command JSON output and signature.
func TestDAOTokenMintJSON(t *testing.T) {
	daoMgr = core.NewDAOManager()
	ledger = core.NewLedger()
	ledger.Mint("admin", 100)
	daoTokenLedger = core.NewDAOTokenLedger(daoMgr, ledger)

	daoMgr.AuthorizeRelayer("admin")
	dao, err := daoMgr.Create("d1", "admin")
	if err != nil {
		t.Fatalf("create dao: %v", err)
	}
	daoMgr.AuthorizeRelayer("addr1")
	if err := daoMgr.Join(dao.ID, "addr1"); err != nil {
		t.Fatalf("join: %v", err)
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	amt := "7"
	msg := []byte(dao.ID + "addr1" + amt)
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

	out, err := execCommand("dao-token", "mint", dao.ID, "admin", "addr1", amt, "--pub", pubHex, "--msg", msgHex, "--sig", sigHex, "--json")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	start := strings.LastIndex(out, "{")
	end := strings.LastIndex(out, "}")
	if start != -1 && end != -1 {
		out = out[start : end+1]
	}
	var resp map[string]string
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["status"] != "minted" {
		t.Fatalf("unexpected status: %v", resp)
	}
	if daoTokenLedger.Balance(dao.ID, "addr1") != 7 {
		t.Fatalf("mint not recorded")
	}
}
