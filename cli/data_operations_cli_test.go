package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDataFeedApplyAndSnapshot(t *testing.T) {
	t.Setenv("SYN_DATA_DIR", t.TempDir())
	resetDataOperationsState()

	manifest := filepath.Join(t.TempDir(), "feed.json")
	if err := os.WriteFile(manifest, []byte(`{"requests":"1200","error_rate":0.01}`), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if _, err := execCommand("data", "feed", "apply", "--feed", "analytics", "--file", manifest); err != nil {
		t.Fatalf("feed apply: %v", err)
	}

	resetDataOperationsState()

	out, err := execCommand("data", "--json", "feed", "snapshot", "analytics")
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	payload := jsonPayload(out)
	var snapshot map[string]string
	if err := json.Unmarshal([]byte(payload), &snapshot); err != nil {
		t.Fatalf("parse snapshot: %v", err)
	}
	if snapshot["requests"] != "1200" {
		t.Fatalf("unexpected requests value: %v", snapshot["requests"])
	}
	if _, err = execCommand("data", "feed", "delete", "analytics", "error_rate"); err != nil {
		t.Fatalf("delete key: %v", err)
	}
	resetDataOperationsState()
	if _, err = execCommand("data", "feed", "get", "analytics", "error_rate"); err == nil {
		t.Fatalf("expected error retrieving deleted key")
	}
}

func TestDataResourceImportAndInfo(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("SYN_DATA_DIR", dataDir)
	resetDataOperationsState()

	assetPath := filepath.Join(dataDir, "asset.txt")
	if err := os.WriteFile(assetPath, []byte("hello world"), 0o600); err != nil {
		t.Fatalf("write asset: %v", err)
	}
	manifest := filepath.Join(dataDir, "resources.json")
	manifestBody := `[
        {"key":"asset","path":"asset.txt","labels":["pii","finance"],"source":"staging"},
        {"key":"inline","data":"memo","labels":["notes"]}
    ]`
	if err := os.WriteFile(manifest, []byte(manifestBody), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if _, err := execCommand("data", "resource", "import", "--manifest", manifest); err != nil {
		t.Fatalf("resource import: %v", err)
	}

	resetDataOperationsState()

	out, err := execCommand("data", "--json", "resource", "info", "asset")
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	payload := jsonPayload(out)
	type info struct {
		Key    string   `json:"key"`
		Size   int      `json:"size"`
		Labels []string `json:"labels"`
		Source string   `json:"source"`
	}
	var detail info
	if err := json.Unmarshal([]byte(payload), &detail); err != nil {
		t.Fatalf("parse info: %v", err)
	}
	if detail.Key != "asset" || detail.Size == 0 || !strings.Contains(strings.Join(detail.Labels, ","), "pii") {
		t.Fatalf("unexpected detail: %#v", detail)
	}

	prunedManifest := filepath.Join(dataDir, "resources_prune.json")
	if err := os.WriteFile(prunedManifest, []byte(`[{"key":"asset"}]`), 0o600); err != nil {
		t.Fatalf("write prune manifest: %v", err)
	}
	if _, err := execCommand("data", "resource", "import", "--manifest", prunedManifest, "--prune"); err != nil {
		t.Fatalf("import prune: %v", err)
	}
	if _, err := execCommand("data", "resource", "get", "inline"); err == nil {
		t.Fatalf("expected inline resource removed")
	}
}
