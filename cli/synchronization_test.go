package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"synnergy/core"
)

func TestSynchronizationCommands(t *testing.T) {
	syncMgr = core.NewSyncManager(ledger)
	syncJSON = false
	t.Cleanup(func() {
		syncMgr = core.NewSyncManager(ledger)
		syncJSON = false
	})

	out, err := executeCLICommand(t, "synchronization", "status")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if !strings.Contains(out, "running: false") {
		t.Fatalf("unexpected status output: %q", out)
	}

	if _, err := executeCLICommand(t, "synchronization", "once"); err == nil || !strings.Contains(err.Error(), "not running") {
		t.Fatalf("expected not running error, got %v", err)
	}

	if _, err := executeCLICommand(t, "synchronization", "start"); err != nil {
		t.Fatalf("start: %v", err)
	}
	if _, err := executeCLICommand(t, "synchronization", "once"); err != nil {
		t.Fatalf("once: %v", err)
	}

	out, err = executeCLICommand(t, "synchronization", "status")
	if err != nil {
		t.Fatalf("status running: %v", err)
	}
	if !strings.Contains(out, "running: true") {
		t.Fatalf("expected running true, got %q", out)
	}

	out, err = executeCLICommand(t, "synchronization", "--json", "status")
	if err != nil {
		t.Fatalf("status json: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}
	if payload["running"] != true {
		t.Fatalf("expected running true in json, got %v", payload["running"])
	}

	if _, err := executeCLICommand(t, "synchronization", "stop"); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
