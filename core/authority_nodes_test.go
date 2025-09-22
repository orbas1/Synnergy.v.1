package core

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"
)

func TestAuthorityNodeRegistry(t *testing.T) {
	reg := NewAuthorityNodeRegistry(nil, NewValidatorManager(MinStake), 0)
	if _, err := reg.Register("addr1", "validator"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if !reg.IsAuthorityNode("addr1") {
		t.Fatalf("expected addr1 to be authority node")
	}
	pub, priv, _ := ed25519.GenerateKey(nil)
	voter := hex.EncodeToString(pub)
	sig := ed25519.Sign(priv, []byte("addr1"))
	if err := reg.Vote(voter, "addr1", sig, pub); err != nil {
		t.Fatalf("vote: %v", err)
	}
	elect := reg.Electorate(1)
	if len(elect) != 1 || elect[0] != "addr1" {
		t.Fatalf("unexpected electorate: %v", elect)
	}
	reg.Deregister("addr1")
	if reg.IsAuthorityNode("addr1") {
		t.Fatalf("deregister failed")
	}
}

func TestAuthorityNodeJSONAndRemoveVote(t *testing.T) {
	reg := NewAuthorityNodeRegistry(nil, NewValidatorManager(MinStake), 0)
	node, err := reg.Register("addr1", "validator")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	pub, priv, _ := ed25519.GenerateKey(nil)
	voter := hex.EncodeToString(pub)
	sig := ed25519.Sign(priv, []byte("addr1"))
	if err := reg.Vote(voter, "addr1", sig, pub); err != nil {
		t.Fatalf("vote: %v", err)
	}
	reg.RemoveVote(voter, "addr1")
	if node.TotalVotes() != 0 {
		t.Fatalf("expected votes cleared")
	}
	if _, err := node.MarshalJSON(); err != nil {
		t.Fatalf("marshal: %v", err)
	}
}
