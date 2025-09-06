package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestDAOCLIFlow verifies basic DAO creation and JSON listing via the CLI.
func TestDAOCLIFlow(t *testing.T) {
	out, err := execCommand("--json", "dao", "create", "testdao", "creator")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	start := strings.LastIndex(out, "{")
	end := strings.LastIndex(out, "}")
	if start != -1 && end != -1 {
		out = out[start : end+1]
	}
	var createRes map[string]interface{}
	if err := json.Unmarshal([]byte(out), &createRes); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	id, _ := createRes["id"].(string)
	if id == "" {
		t.Fatalf("expected DAO id output")
	}

	if _, err := execCommand("--json", "dao", "join", id, "member1"); err != nil {
		t.Fatalf("join: %v", err)
	}

	out, err = execCommand("--json", "dao", "info", id)
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	start = strings.LastIndex(out, "{")
	end = strings.LastIndex(out, "}")
	if start != -1 && end != -1 {
		out = out[start : end+1]
	}
	var info map[string]interface{}
	if err := json.Unmarshal([]byte(out), &info); err != nil {
		t.Fatalf("decode info: %v", err)
	}
	if info["name"] != "testdao" {
		t.Fatalf("unexpected info name: %v", info)
	}

	out, err = execCommand("--json", "dao", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	s := strings.Index(out, "[")
	e := strings.LastIndex(out, "]")
	if s != -1 && e != -1 && e >= s {
		out = out[s : e+1]
	}
	var list []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("expected at least one DAO")
	}
}
