package cli

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"synnergy/core"
)

func createWalletFile(t *testing.T, password string) (*core.Wallet, string) {
	t.Helper()
	w, path := newMemoryWallet(t, password)
	return w, path
}

func TestSyn3800Lifecycle(t *testing.T) {
	useMemoryWalletLoader(t)
	setStage73StatePath(filepath.Join(t.TempDir(), "stage73.json"))
	resetStage73LoadedForTests()
	authorizer, path := createWalletFile(t, "pass")
	out, err := execCommand("syn3800", "create", "bob", "research", "100", "--authorizer", path+":pass")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if firstNonGasLine(out) != "1" {
		t.Fatalf("expected id 1, got %s", out)
	}
	if _, err := execCommand("syn3800", "release", "1", "40", "phase1", "--wallet", path, "--password", "pass"); err != nil {
		t.Fatalf("release: %v", err)
	}
	data, err := execCommand("syn3800", "get", "1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	var grant struct {
		Released uint64
		Status   string
	}
	if err := json.Unmarshal([]byte(jsonPayload(data)), &grant); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if grant.Released != 40 || grant.Status != string(core.GrantStatusActive) {
		t.Fatalf("unexpected grant: %+v", grant)
	}
	second, secondPath := createWalletFile(t, "pw2")
	if _, err := execCommand("syn3800", "authorize", "1", "--wallet", secondPath, "--password", "pw2"); err != nil {
		t.Fatalf("authorize: %v", err)
	}
	if _, err := execCommand("syn3800", "release", "1", "60", "--wallet", secondPath, "--password", "pw2"); err != nil {
		t.Fatalf("release second: %v", err)
	}
	auditJSON, err := execCommand("syn3800", "audit", "1")
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	var events []core.GrantEvent
	if err := json.Unmarshal([]byte(jsonPayload(auditJSON)), &events); err != nil {
		t.Fatalf("audit unmarshal: %v", err)
	}
	if len(events) < 3 {
		t.Fatalf("expected audit events, got %d", len(events))
	}
	// ensure CLI uses uppercase symbol for tokens by verifying release signer address recorded
	if events[len(events)-1].Amount != 60 {
		t.Fatalf("expected final disbursement of 60, got %+v", events[len(events)-1])
	}
	// use second wallet variable to silence unused warning
	if second.Address == "" || authorizer.Address == "" {
		t.Fatalf("wallet addresses missing")
	}
	statusJSON, err := execCommand("syn3800", "status")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	var status struct {
		Total     int
		Pending   int
		Active    int
		Completed int
	}
	if err := json.Unmarshal([]byte(jsonPayload(statusJSON)), &status); err != nil {
		t.Fatalf("unmarshal status: %v", err)
	}
	if status.Total != 1 || status.Completed != 1 {
		t.Fatalf("unexpected telemetry output: %+v", status)
	}
}

func TestSyn3800Validation(t *testing.T) {
	useMemoryWalletLoader(t)
	setStage73StatePath(filepath.Join(t.TempDir(), "stage73.json"))
	resetStage73LoadedForTests()
	if _, err := execCommand("syn3800", "create", "", "name", "10"); err == nil {
		t.Fatal("expected error for missing beneficiary")
	} else if !strings.Contains(err.Error(), "beneficiary") {
		t.Fatalf("unexpected beneficiary error: %v", err)
	}
	if _, err := execCommand("syn3800", "create", "bob", "", "10"); err == nil {
		t.Fatal("expected error for missing name")
	} else if !strings.Contains(err.Error(), "name required") {
		t.Fatalf("unexpected name error: %v", err)
	}
	if _, err := execCommand("syn3800", "create", "bob", "research", "0"); err == nil {
		t.Fatal("expected error for zero amount")
	} else if !strings.Contains(err.Error(), "invalid amount") {
		t.Fatalf("unexpected amount error: %v", err)
	}
	_, path := createWalletFile(t, "pass")
	if _, err := execCommand("syn3800", "create", "alice", "grant", "50", "--authorizer", path+":pass"); err != nil {
		t.Fatalf("create valid: %v", err)
	}
	if _, err := execCommand("syn3800", "release", "1", "10", "--wallet", path, "--password", "wrong"); err == nil {
		t.Fatal("expected error for wrong password")
	} else if !strings.Contains(err.Error(), "authentication failed") {
		t.Fatalf("unexpected wrong password error: %v", err)
	}
	listOut, err := execCommand("syn3800", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	var grants []struct{ ID int }
	if err := json.Unmarshal([]byte(jsonPayload(listOut)), &grants); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if len(grants) != 1 || grants[0].ID != 1 {
		t.Fatalf("unexpected list data: %s", listOut)
	}
	auditOut, err := execCommand("syn3800", "audit", "1")
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	if auditOut == "" {
		t.Fatal("expected audit log output")
	}
	// ensure release fails without wallet flags
	if _, err := execCommand("syn3800", "release", "1", "10"); err == nil {
		t.Fatal("expected error for missing wallet flags")
	} else if !strings.Contains(err.Error(), "wallet and password required") {
		t.Fatalf("unexpected missing wallet error: %v", err)
	}
}
