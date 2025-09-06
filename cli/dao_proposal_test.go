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

// TestDAOProposalCreateJSON verifies JSON output and signature check for proposal creation.
func TestDAOProposalCreateJSON(t *testing.T) {
	daoMgr = core.NewDAOManager()
	proposalMgr = core.NewProposalManager()
	dao := daoMgr.Create("dao1", "creator")

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	desc := "proposal"
	msg := []byte(dao.ID + "creator" + desc)
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
	rootCmd.SetArgs([]string{"dao-proposal", "create", dao.ID, "creator", desc, "--pub", pubHex, "--msg", msgHex, "--sig", sigHex, "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	var resp map[string]string
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["id"] == "" {
		t.Fatalf("missing id: %v", resp)
	}
	if len(proposalMgr.List()) != 1 {
		t.Fatalf("proposal not created")
	}
}
