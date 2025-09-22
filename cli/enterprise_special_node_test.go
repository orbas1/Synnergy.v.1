package cli

import (
	"testing"
)

func TestEnterpriseSpecialCommandRegistered(t *testing.T) {
	root := RootCmd()
	if root == nil {
		t.Fatalf("root command not initialised")
	}
	found := false
	for _, cmd := range root.Commands() {
		if cmd.Name() == "enterprise-special" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("enterprise-special command not registered")
	}
}

func TestEnterpriseSpecialAttachCommand(t *testing.T) {
	cmd := newEnterpriseAttachCmd()
	if err := cmd.Flags().Set("id", "cli-attach-test"); err != nil {
		t.Fatalf("set id: %v", err)
	}
	if err := cmd.Flags().Set("role", "analytics"); err != nil {
		t.Fatalf("set role: %v", err)
	}
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("run attach: %v", err)
	}
	if !enterpriseSpecialNode.DetachPlugin("cli-attach-test") {
		t.Fatalf("expected plugin to be registered")
	}
}

func TestEnterpriseSpecialBroadcastCommand(t *testing.T) {
	attach := newEnterpriseAttachCmd()
	if err := attach.Flags().Set("id", "cli-broadcast"); err != nil {
		t.Fatalf("set id: %v", err)
	}
	if err := attach.Flags().Set("seed-balance", "250"); err != nil {
		t.Fatalf("set seed: %v", err)
	}
	if err := attach.RunE(attach, nil); err != nil {
		t.Fatalf("attach: %v", err)
	}

	broadcast := newEnterpriseBroadcastCmd()
	if err := broadcast.Flags().Set("from", "cli-broadcast"); err != nil {
		t.Fatalf("from: %v", err)
	}
	if err := broadcast.Flags().Set("to", "cli-dest"); err != nil {
		t.Fatalf("to: %v", err)
	}
	if err := broadcast.Flags().Set("amount", "10"); err != nil {
		t.Fatalf("amount: %v", err)
	}
	if err := broadcast.Flags().Set("fee", "1"); err != nil {
		t.Fatalf("fee: %v", err)
	}
	if err := broadcast.Flags().Set("nonce", "1"); err != nil {
		t.Fatalf("nonce: %v", err)
	}
	if err := broadcast.RunE(broadcast, nil); err != nil {
		t.Fatalf("broadcast: %v", err)
	}

	// Ensure the transaction was queued for at least one plugin.
	snap := enterpriseSpecialNode.Snapshot()
	queued := false
	for _, plugin := range snap.Plugins {
		if plugin.ID == "cli-broadcast" && plugin.Metrics.MempoolSize > 0 {
			queued = true
			break
		}
	}
	if !queued {
		t.Fatalf("expected broadcast to queue transaction")
	}

	enterpriseSpecialNode.DetachPlugin("cli-broadcast")
}
