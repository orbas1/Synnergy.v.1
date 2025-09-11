package cli

import (
	"bytes"
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

// TestDAOQuadraticWeightJSON verifies JSON output for weight calculation.
func TestDAOQuadraticWeightJSON(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"dao-qv", "weight", "9", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	var resp map[string]uint64
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	expected := core.QuadraticWeight(9)
	if resp["weight"] != expected {
		t.Fatalf("unexpected weight: %v", resp)
	}
}

// TestDAOQuadraticVoteRequiresMembership ensures votes from non-members are rejected.
func TestDAOQuadraticVoteRequiresMembership(t *testing.T) {
	daoMgr = core.NewDAOManager()
	daoMgr.AuthorizeRelayer("c")
	proposalMgr = core.NewProposalManager()
	dao, err := daoMgr.Create("dao1", "c")
	if err != nil {
		t.Fatalf("create dao: %v", err)
	}
	p, err := proposalMgr.CreateProposal(dao, "c", "desc")
	if err != nil {
		t.Fatalf("create proposal: %v", err)
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	msg := []byte(p.ID + "v" + "4" + "yes")
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
	rootCmd.SetArgs([]string{"dao-qv", "vote", p.ID, "v", "4", "yes", "--pub", pubHex, "--msg", msgHex, "--sig", sigHex})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "not a dao member") {
		t.Fatalf("expected membership error, got %q", buf.String())
	}
}
