package tests

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"synnergy/core"
)

func runStage73CLI(t *testing.T, statePath string, args ...string) string {
	t.Helper()
	absRepo, err := filepath.Abs(repoPath())
	if err != nil {
		t.Fatalf("abs repo path: %v", err)
	}
	mainPath := filepath.Join(absRepo, "cmd", "synnergy", "main.go")
	base := []string{"run", mainPath, "--stage73-state", statePath}
	base = append(base, args...)
	cmd := exec.Command("go", base...)
	cmd.Dir = absRepo
	cmd.Env = append(os.Environ(), fmt.Sprintf("SYN_CONFIG=%s", filepath.Join(absRepo, "configs", "test.yaml")))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command %v failed: %v\n%s", args, err, string(out))
	}
	return string(out)
}

func TestStage73EndToEndWorkflow(t *testing.T) {
	dir := t.TempDir()
	statePath := filepath.Join(dir, "stage73_state.json")
	walletPath := filepath.Join(dir, "controller.wallet")

	runStage73CLI(t, statePath, "wallet", "new", "--out", walletPath, "--password", "pass")
	runStage73CLI(t, statePath, "syn3700", "init", "--name", "Institutional", "--symbol", "IDX", "--controller", walletPath+":pass")
	runStage73CLI(t, statePath, "syn3700", "add", "AAA", "--weight", "0.5", "--drift", "0.1", "--wallet", walletPath, "--password", "pass")
	runStage73CLI(t, statePath, "syn3800", "create", "beneficiary", "Education", "100", "--authorizer", walletPath+":pass")
	runStage73CLI(t, statePath, "syn3800", "release", "1", "40", "phase-one", "--wallet", walletPath, "--password", "pass")
	runStage73CLI(t, statePath, "syn3900", "register", "recipient", "Healthcare", "60", "--approver", walletPath+":pass")
	runStage73CLI(t, statePath, "syn3900", "approve", "1", "--wallet", walletPath, "--password", "pass")
	expiry := time.Now().Add(2 * time.Hour).Unix()
	runStage73CLI(t, statePath, "syn4700", "create", "--id", "AGR-1", "--name", "Agreement", "--symbol", "AGR", "--doctype", "contract", "--hash", "hash", "--owner", "owner-wallet", "--expiry", fmt.Sprint(expiry), "--supply", "10", "--party", "alice", "--party", "bob")
	runStage73CLI(t, statePath, "syn4700", "sign", "AGR-1", "alice", "sig-alice")
	runStage73CLI(t, statePath, "syn4200_token", "donate", "HELP", "--from", "donor", "--amt", "25", "--purpose", "relief")
	runStage73CLI(t, statePath, "syn500", "create", "--name", "Service Credit", "--symbol", "UTL", "--owner", "owner-wallet", "--dec", "4", "--supply", "1000")
	runStage73CLI(t, statePath, "syn500", "grant", "consumer", "--tier", "1", "--max", "5", "--window", "30m")
	runStage73CLI(t, statePath, "syn500", "use", "consumer")

	store := core.NewStage73Store(statePath)
	if err := store.Load(); err != nil {
		t.Fatalf("load snapshot: %v", err)
	}
	snap := store.Snapshot()
	if snap.Index == nil || len(snap.Index.Components) == 0 {
		t.Fatalf("expected index snapshot, got %+v", snap.Index)
	}
	if snap.Grants.Summary.Total != 1 || snap.Grants.Records[0].Released == 0 {
		t.Fatalf("unexpected grant snapshot: %+v", snap.Grants)
	}
	if snap.Benefits.Summary.Total != 1 || snap.Benefits.Summary.Approved == 0 {
		t.Fatalf("unexpected benefit snapshot: %+v", snap.Benefits.Summary)
	}
	if len(snap.Charity) == 0 {
		t.Fatalf("expected charity donation recorded")
	}
	if snap.Utility == nil || snap.Utility.Telemetry.Grants == 0 {
		t.Fatalf("unexpected utility telemetry: %+v", snap.Utility)
	}

	orchestrator := core.NewStage73Orchestrator(store, core.NewSynnergyConsensus())
	digest, err := orchestrator.SnapshotDigest()
	if err != nil {
		t.Fatalf("orchestrator digest: %v", err)
	}
	if digest == "" {
		t.Fatal("expected non-empty digest")
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}
	var canonicalSnap core.Stage73Snapshot
	if err := json.Unmarshal(data, &canonicalSnap); err != nil {
		t.Fatalf("unmarshal snapshot: %v", err)
	}
	canonical, err := json.Marshal(canonicalSnap)
	if err != nil {
		t.Fatalf("marshal canonical snapshot: %v", err)
	}
	canonicalHash := sha256.Sum256(canonical)
	if fmt.Sprintf("%x", canonicalHash) != digest {
		t.Fatalf("digest mismatch: orchestrator=%s api=%x", digest, canonicalHash)
	}
}
