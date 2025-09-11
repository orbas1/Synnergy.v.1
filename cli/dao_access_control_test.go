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

// TestDAOMemberAddJSON verifies JSON output and signature verification for member addition.
func TestDAOMemberAddJSON(t *testing.T) {
	daoMgr = core.NewDAOManager()
	daoMgr.AuthorizeRelayer("creator")
	dao, err := daoMgr.Create("testdao", "creator")
	if err != nil {
		t.Fatalf("failed to create dao: %v", err)
	}

	// generate key and signature over daoID+addr+role
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	msg := []byte(dao.ID + "member1" + "member")
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
	rootCmd.SetArgs([]string{"dao-members", "add", dao.ID, "member1", "member", "--pub", pubHex, "--msg", msgHex, "--sig", sigHex, "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	var resp map[string]string
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["status"] != "member added" {
		t.Fatalf("unexpected status: %v", resp)
	}

	if role, ok := dao.MemberRole("member1"); !ok || role != "member" {
		t.Fatalf("member not added: %v %v", role, ok)
	}
}

// TestDAOMemberUpdateRole ensures only admins can update roles and JSON output is emitted.
func TestDAOMemberUpdateRole(t *testing.T) {
	daoMgr = core.NewDAOManager()
	daoMgr.AuthorizeRelayer("creator")
	dao, err := daoMgr.Create("testdao", "creator")
	if err != nil {
		t.Fatalf("create dao: %v", err)
	}
	if err := dao.AddMember("member1", "member"); err != nil {
		t.Fatalf("add member: %v", err)
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	msg := []byte(dao.ID + "creator" + "member1" + "admin")
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
	rootCmd.SetArgs([]string{"dao-members", "update", dao.ID, "creator", "member1", "admin", "--pub", pubHex, "--msg", msgHex, "--sig", sigHex, "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	var resp map[string]string
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["status"] != "role updated" {
		t.Fatalf("unexpected status: %v", resp)
	}
	if role, _ := dao.MemberRole("member1"); role != "admin" {
		t.Fatalf("role not updated: %v", role)
	}
}

// TestDAOMemberUpdateRoleUnauthorized rejects updates from non-admins.
func TestDAOMemberUpdateRoleUnauthorized(t *testing.T) {
	daoMgr = core.NewDAOManager()
	daoMgr.AuthorizeRelayer("creator")
	daoMgr.AuthorizeRelayer("member1")
	dao, err := daoMgr.Create("testdao", "creator")
	if err != nil {
		t.Fatalf("create dao: %v", err)
	}
	if err := dao.AddMember("member1", "member"); err != nil {
		t.Fatalf("add member: %v", err)
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	msg := []byte(dao.ID + "member1" + "creator" + "member")
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
	rootCmd.SetArgs([]string{"dao-members", "update", dao.ID, "member1", "creator", "member", "--pub", pubHex, "--msg", msgHex, "--sig", sigHex})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("unauthorized")) {
		t.Fatalf("expected unauthorized error, got %q", buf.String())
	}
}
