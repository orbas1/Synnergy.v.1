package core

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

func TestAuthorityApplication(t *testing.T) {
	reg := NewAuthorityNodeRegistry(nil, NewValidatorManager(MinStake), 0)
	mgr := NewAuthorityApplicationManager(reg, time.Hour)
	id := mgr.Submit("cand1", "validator", "test")

	pub, priv, _ := ed25519.GenerateKey(nil)
	voterAddr := hex.EncodeToString(pub)
	msg := fmt.Sprintf("%s:%t", id, true)
	sig := ed25519.Sign(priv, []byte(msg))

	if err := mgr.Vote(voterAddr, id, true, sig, pub); err != nil {
		t.Fatalf("vote: %v", err)
	}
	if err := mgr.Finalize(id); err != nil {
		t.Fatalf("finalize: %v", err)
	}
	if !reg.IsAuthorityNode("cand1") {
		t.Fatalf("candidate not registered after finalise")
	}
}

func TestAuthorityApplicationJSON(t *testing.T) {
	reg := NewAuthorityNodeRegistry(nil, NewValidatorManager(MinStake), 0)
	mgr := NewAuthorityApplicationManager(reg, time.Hour)
	id := mgr.Submit("cand1", "validator", "test")
	if id == "" {
		t.Fatalf("empty application id")
	}
	app, err := mgr.Get(id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if _, err := app.MarshalJSON(); err != nil {
		t.Fatalf("marshal: %v", err)
	}
}
