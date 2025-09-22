package tokens

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"testing"
)

func TestCarbonRegistryLifecycle(t *testing.T) {
	reg := NewCarbonRegistry()
	proj := reg.Register("alice", "Mangrove", 1_000)

	if err := reg.Issue(proj.ID, "alice", 400); err != nil {
		t.Fatalf("issue: %v", err)
	}
	if err := reg.Transfer(proj.ID, "alice", "bob", 150); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if err := reg.Retire(proj.ID, "bob", 50); err != nil {
		t.Fatalf("retire: %v", err)
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	reg.RegisterVerifier("auditor", pub)

	payload := fmt.Sprintf("%s|auditor|V1|ok", proj.ID)
	sig := ed25519.Sign(priv, []byte(payload))
	if err := reg.AddSignedVerification(proj.ID, "auditor", "V1", "ok", sig); err != nil {
		t.Fatalf("add signed verification: %v", err)
	}

	if err := reg.AddSignedVerification(proj.ID, "auditor", "V1b", "ok", []byte("bad")); err == nil {
		t.Fatal("expected invalid signature error")
	}

	snapshot, err := reg.Snapshot(proj.ID)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if snapshot.HolderBalances["alice"] != 250 {
		t.Fatalf("unexpected balance snapshot: %+v", snapshot)
	}

	events, ok := reg.ProjectTimeline(proj.ID, 10)
	if !ok || len(events) < 3 {
		t.Fatalf("expected events, got %+v", events)
	}
	if events[len(events)-1].Type != "verification" {
		t.Fatalf("expected latest event to be verification, got %+v", events[len(events)-1])
	}
}
