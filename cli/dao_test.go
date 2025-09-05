package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestDAOCLIFlow verifies basic DAO creation and JSON listing via the CLI.
func TestDAOCLIFlow(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"dao", "create", "testdao", "creator"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	id := lines[len(lines)-1]
	if id == "" {
		t.Fatalf("expected DAO id output")
	}

	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"dao", "join", id, "member1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("join: %v", err)
	}

	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"dao", "info", id, "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("info: %v", err)
	}
	var info map[string]interface{}
	if err := json.NewDecoder(buf).Decode(&info); err != nil {
		t.Fatalf("decode info: %v", err)
	}
	if info["Name"] != "testdao" {
		t.Fatalf("unexpected info name: %v", info)
	}

	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"dao", "list", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}
	var list []map[string]interface{}
	if err := json.NewDecoder(buf).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("expected at least one DAO")
	}
}
