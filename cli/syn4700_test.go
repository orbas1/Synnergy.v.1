package cli

import (
	"encoding/json"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"synnergy/core"
)

func TestSyn4700Lifecycle(t *testing.T) {
	setStage73StatePath(filepath.Join(t.TempDir(), "stage73.json"))
	resetStage73LoadedForTests()
	expiry := time.Now().Add(time.Hour).Unix()
	if _, err := execCommand("syn4700", "create", "--id", "t1", "--name", "Agreement", "--symbol", "AGR", "--doctype", "contract", "--hash", "h", "--owner", "alice", "--expiry", strconv.FormatInt(expiry, 10), "--supply", "100", "--party", "alice", "--party", "bob"); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := execCommand("syn4700", "sign", "t1", "alice", "sig"); err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := execCommand("syn4700", "status", "t1", "active"); err != nil {
		t.Fatalf("status: %v", err)
	}
	if _, err := execCommand("syn4700", "dispute", "t1", "breach", "resolved"); err != nil {
		t.Fatalf("dispute: %v", err)
	}
	out, err := execCommand("syn4700", "info", "t1")
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	var tok struct {
		Status     string
		Signatures map[string]string
	}
	if err := json.Unmarshal([]byte(out), &tok); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if tok.Status != string(core.LegalTokenStatusDisputed) {
		t.Fatalf("expected disputed status, got %s", tok.Status)
	}
	if tok.Signatures["alice"] != "sig" {
		t.Fatalf("signature missing")
	}
}

func TestSyn4700Validation(t *testing.T) {
	setStage73StatePath(filepath.Join(t.TempDir(), "stage73.json"))
	resetStage73LoadedForTests()
	expiry := time.Now().Add(time.Hour).Unix()
	if _, err := execCommand("syn4700", "create", "--id", "t1", "--name", "A", "--symbol", "A", "--doctype", "d", "--hash", "h", "--owner", "o", "--expiry", strconv.FormatInt(expiry, 10), "--supply", "0", "--party", "p1"); err == nil {
		t.Fatal("expected error for zero supply")
	}
	if _, err := execCommand("syn4700", "create", "--id", "t1", "--name", "A", "--symbol", "A", "--doctype", "d", "--hash", "h", "--owner", "o", "--expiry", strconv.FormatInt(expiry, 10), "--supply", "10", "--party", "p1"); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := execCommand("syn4700", "sign", "t1", "unknown", "sig"); err == nil {
		t.Fatal("expected error for unknown party")
	}
}
