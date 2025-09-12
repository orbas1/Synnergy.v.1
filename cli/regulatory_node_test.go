package cli

import (
	"os"
	"testing"

	"synnergy/core"
)

func TestRegNodeApproveRequiresSignedWallet(t *testing.T) {
	regManager = core.NewRegulatoryManager()
	regManager.AddRegulation(core.Regulation{ID: "r1", MaxAmount: 10})
	regNode = core.NewRegulatoryNode("regnode1", regManager)

	w, err := core.NewWallet()
	if err != nil {
		t.Fatalf("wallet: %v", err)
	}
	path := tempFile(t)
	if err := w.Save(path, "pw"); err != nil {
		t.Fatalf("save: %v", err)
	}
	cmd := RootCmd()
	cmd.SetArgs([]string{"regnode", "approve", w.Address, "5", "--wallet", path, "--password", "pw"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("approve failed: %v", err)
	}
	if logs := regNode.Logs(w.Address); len(logs) != 0 {
		t.Fatalf("expected no logs, got %v", logs)
	}

	cmd.SetArgs([]string{"regnode", "approve", w.Address, "20", "--wallet", path, "--password", "pw"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("approve failed: %v", err)
	}
	if logs := regNode.Logs(w.Address); len(logs) != 1 {
		t.Fatalf("expected rejection log, got %v", logs)
	}
}

func tempFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("", "wallet")
	if err != nil {
		t.Fatalf("tempfile: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestRegNodeFlagRequiresReason(t *testing.T) {
	regManager = core.NewRegulatoryManager()
	regNode = core.NewRegulatoryNode("regnode1", regManager)
	cmd := RootCmd()
	cmd.SetArgs([]string{"regnode", "flag", "eve", ""})
	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error for empty reason")
	}
	if logs := regNode.Logs("eve"); len(logs) != 0 {
		t.Fatalf("expected no logs, got %v", logs)
	}
}
