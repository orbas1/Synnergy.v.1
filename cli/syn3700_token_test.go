package cli

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"synnergy/core"
)

func createIndexWallet(t *testing.T, password string) (*core.Wallet, string) {
	t.Helper()
	wallet, path := newMemoryWallet(t, password)
	return wallet, path
}

func TestSyn3700CLIWorkflow(t *testing.T) {
	useMemoryWalletLoader(t)
	setStage73StatePath(filepath.Join(t.TempDir(), "stage73.json"))
	resetStage73LoadedForTests()
	syn3700 = nil
	controller, path := createIndexWallet(t, "pass")
	if out, err := execCommand("syn3700", "init", "--name", "Institutional", "--symbol", "IDX", "--controller", path+":pass"); err != nil {
		t.Fatalf("init: %v (output: %s)", err, out)
	}
	if out, err := execCommand("syn3700", "add", "AAA", "--weight", "0.5", "--drift", "0.1", "--wallet", path, "--password", "pass"); err != nil {
		t.Fatalf("add AAA: %v (output: %s)", err, out)
	}
	if out, err := execCommand("syn3700", "add", "BBB", "--weight", "1.5", "--drift", "0.2", "--wallet", path, "--password", "pass"); err != nil {
		t.Fatalf("add BBB: %v (output: %s)", err, out)
	}
	snapJSON, err := execCommand("syn3700", "snapshot")
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	var snap struct {
		Symbol     string           `json:"symbol"`
		Components []core.Component `json:"components"`
	}
	if err := json.Unmarshal([]byte(jsonPayload(snapJSON)), &snap); err != nil {
		t.Fatalf("unmarshal snapshot: %v", err)
	}
	if snap.Symbol != "IDX" || len(snap.Components) != 2 {
		t.Fatalf("unexpected snapshot: %+v", snap)
	}
	statusJSON, err := execCommand("syn3700", "status")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	var telemetry struct {
		ComponentCount  int `json:"component_count"`
		ControllerCount int `json:"controller_count"`
	}
	if err := json.Unmarshal([]byte(jsonPayload(statusJSON)), &telemetry); err != nil {
		t.Fatalf("unmarshal status: %v", err)
	}
	if telemetry.ComponentCount != 2 || telemetry.ControllerCount != 1 {
		t.Fatalf("unexpected telemetry: %+v", telemetry)
	}
	controllersJSON, err := execCommand("syn3700", "controllers")
	if err != nil {
		t.Fatalf("controllers: %v", err)
	}
	var controllers []string
	if err := json.Unmarshal([]byte(jsonPayload(controllersJSON)), &controllers); err != nil {
		t.Fatalf("unmarshal controllers: %v", err)
	}
	if len(controllers) != 1 || controllers[0] != controller.Address {
		t.Fatalf("unexpected controllers: %v", controllers)
	}
	valueJSON, err := execCommand("syn3700", "value", "AAA:10", "BBB:20")
	if err != nil {
		t.Fatalf("value: %v", err)
	}
	var report struct {
		Value float64 `json:"value"`
	}
	if err := json.Unmarshal([]byte(jsonPayload(valueJSON)), &report); err != nil {
		t.Fatalf("unmarshal report: %v", err)
	}
	if report.Value != 35 {
		t.Fatalf("expected value 35, got %f", report.Value)
	}
	rebalanceJSON, err := execCommand("syn3700", "rebalance", "--wallet", path, "--password", "pass")
	if err != nil {
		t.Fatalf("rebalance: %v", err)
	}
	var changes map[string][2]float64
	if err := json.Unmarshal([]byte(jsonPayload(rebalanceJSON)), &changes); err != nil {
		t.Fatalf("unmarshal rebalance: %v", err)
	}
	if len(changes) != 2 {
		t.Fatalf("expected 2 rebalance entries")
	}
	auditJSON, err := execCommand("syn3700", "audit")
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	if auditJSON == "" {
		t.Fatalf("expected audit output")
	}
}

func TestSyn3700CLIValidation(t *testing.T) {
	useMemoryWalletLoader(t)
	setStage73StatePath(filepath.Join(t.TempDir(), "stage73.json"))
	resetStage73LoadedForTests()
	syn3700 = nil
	_, path := createIndexWallet(t, "pass")
	if _, err := execCommand("syn3700", "add", "AAA", "--weight", "1", "--wallet", path, "--password", "pass"); err == nil {
		t.Fatal("expected error when not initialised")
	} else if !strings.Contains(err.Error(), "token not initialised") {
		t.Fatalf("unexpected init error: %v", err)
	}
	if _, err := execCommand("syn3700", "init", "--name", "Inst", "--symbol", "IDX", "--controller", path+":pass"); err != nil {
		t.Fatalf("init: %v", err)
	}
	if _, err := execCommand("syn3700", "add", "AAA", "--weight", "0", "--drift", "0.1", "--wallet", path, "--password", "pass"); err == nil {
		t.Fatal("expected weight validation error")
	} else if !strings.Contains(err.Error(), "invalid weight") {
		t.Fatalf("unexpected weight error: %v", err)
	}
	if _, err := execCommand("syn3700", "add", "AAA", "--weight", "1", "--drift", "1.1", "--wallet", path, "--password", "pass"); err == nil {
		t.Fatal("expected drift validation error")
	} else if !strings.Contains(err.Error(), "component drift") {
		t.Fatalf("unexpected drift error: %v", err)
	}
	if _, err := execCommand("syn3700", "remove", "AAA", "--wallet", path, "--password", "pass"); err == nil {
		t.Fatal("expected removal validation error")
	} else if !strings.Contains(err.Error(), "component not found") {
		t.Fatalf("unexpected removal error: %v", err)
	}
	if _, err := execCommand("syn3700", "rebalance"); err == nil {
		t.Fatal("expected error for missing wallet flags")
	} else if !strings.Contains(err.Error(), "required flag") {
		t.Fatalf("unexpected rebalance error: %v", err)
	}
}
