package cli

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"testing"
	"time"

	"synnergy/core"
)

func TestParseMetadataDeterministic(t *testing.T) {
	m := parseMetadata([]string{"unit=alpha", "priority=high", "unit=override"})
	if len(m) != 2 {
		t.Fatalf("expected deduplicated map, got %v", m)
	}
	if m["unit"] != "override" {
		t.Fatalf("expected last value to win, got %v", m["unit"])
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	if keys[0] != "priority" {
		t.Fatalf("expected deterministic ordering, got %v", keys)
	}
}

func TestWarfareNodeManualCommandFlow(t *testing.T) {
	base := core.NewNode("cli", "addr", core.NewLedger())
	warfareNode = core.NewWarfareNode(base)
	defer func() { warfareNode = nil }()
	cred, err := warfareNode.IssueCommander("cli-user")
	if err != nil {
		t.Fatalf("issue commander: %v", err)
	}
	priv, err := hex.DecodeString(cred.PrivateKey)
	if err != nil {
		t.Fatalf("decode key: %v", err)
	}
	req := core.CommandRequest{
		Commander: cred.ID,
		Command:   "deploy",
		Timestamp: time.Now().UTC(),
		Metadata:  map[string]string{"cli": "test"},
	}
	req.Signature = ed25519.Sign(ed25519.PrivateKey(priv), req.CanonicalPayload())
	record, err := warfareNode.ExecuteSecureCommand(context.Background(), req)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}
	if !record.Accepted {
		t.Fatalf("expected command accepted")
	}

	if _, err := warfareNode.RecordLogistics(core.LogisticsUpdate{AssetID: "asset", Location: "loc", Status: "ok"}); err != nil {
		t.Fatalf("record logistics: %v", err)
	}
	if len(warfareNode.Logistics()) != 1 {
		t.Fatalf("expected logistics entry")
	}
}
