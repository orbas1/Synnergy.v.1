package cli

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"synnergy/core"
)

func createClaimWallet(t *testing.T, password string) (*core.Wallet, string) {
	t.Helper()
	w, path := newMemoryWallet(t, password)
	return w, path
}

func TestSyn3900Lifecycle(t *testing.T) {
	useMemoryWalletLoader(t)
	setStage73StatePath(filepath.Join(t.TempDir(), "stage73.json"))
	resetStage73LoadedForTests()
	recipient, path := createClaimWallet(t, "pass")
	out, err := execCommand("syn3900", "register", recipient.Address, "housing", "200", "--approver", path+":pass")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if firstNonGasLine(out) != "1" {
		t.Fatalf("expected id 1, got %s", out)
	}
	if _, err := execCommand("syn3900", "claim", "1", "--wallet", path, "--password", "pass"); err != nil {
		t.Fatalf("claim: %v", err)
	}
	data, err := execCommand("syn3900", "get", "1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	var benefit struct {
		Claimed bool
		Status  string
	}
	if err := json.Unmarshal([]byte(jsonPayload(data)), &benefit); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !benefit.Claimed || benefit.Status != string(core.BenefitStatusClaimed) {
		t.Fatalf("unexpected benefit: %+v", benefit)
	}
	approver, path2 := createClaimWallet(t, "pw2")
	if _, err := execCommand("syn3900", "approve", "1", "--wallet", path2, "--password", "pw2"); err != nil {
		t.Fatalf("approve: %v", err)
	}
	// ensure approving after claim is idempotent and wallet is used
	if approver.Address == "" {
		t.Fatalf("approver wallet missing")
	}
	statusJSON, err := execCommand("syn3900", "status")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	var tele struct {
		Total    int
		Pending  int
		Approved int
		Claimed  int
	}
	if err := json.Unmarshal([]byte(jsonPayload(statusJSON)), &tele); err != nil {
		t.Fatalf("unmarshal telemetry: %v", err)
	}
	if tele.Total != 1 || tele.Approved == 0 {
		t.Fatalf("unexpected telemetry output: %+v", tele)
	}
	listJSON, err := execCommand("syn3900", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	listJSON = jsonPayload(listJSON)
	if !strings.Contains(listJSON, "housing") {
		t.Fatalf("expected list output, got %s", listJSON)
	}
}

func TestSyn3900Validation(t *testing.T) {
	useMemoryWalletLoader(t)
	setStage73StatePath(filepath.Join(t.TempDir(), "stage73.json"))
	resetStage73LoadedForTests()
	if _, err := execCommand("syn3900", "register", "", "program", "10"); err == nil {
		t.Fatal("expected error for missing recipient")
	} else if !strings.Contains(err.Error(), "recipient") {
		t.Fatalf("unexpected recipient error: %v", err)
	}
	if _, err := execCommand("syn3900", "register", "alice", "", "10"); err == nil {
		t.Fatal("expected error for missing program")
	} else if !strings.Contains(err.Error(), "program") {
		t.Fatalf("unexpected program error: %v", err)
	}
	if _, err := execCommand("syn3900", "register", "alice", "program", "0"); err == nil {
		t.Fatal("expected error for amount")
	} else if !strings.Contains(err.Error(), "invalid amount") {
		t.Fatalf("unexpected amount error: %v", err)
	}
	recipient, path := createClaimWallet(t, "pass")
	if _, err := execCommand("syn3900", "register", recipient.Address, "benefit", "50"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, err := execCommand("syn3900", "claim", "1", "--wallet", path, "--password", "wrong"); err == nil {
		t.Fatal("expected error for wrong password")
	} else if !strings.Contains(err.Error(), "authentication failed") {
		t.Fatalf("unexpected wrong password error: %v", err)
	}
	if _, err := execCommand("syn3900", "claim", "1"); err == nil {
		t.Fatal("expected error for missing wallet flags")
	} else if !strings.Contains(err.Error(), "wallet and password required") {
		t.Fatalf("unexpected missing wallet error: %v", err)
	}
}
