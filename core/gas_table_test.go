package core

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// gasListPath returns the absolute path to gas_table_list.md used by gas table utilities.
func gasListPath(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("cannot determine caller path")
	}
	coreDir := filepath.Dir(filename)
	rootDir := filepath.Dir(coreDir)
	return filepath.Join(rootDir, "docs", "reference", "gas_table_list.md")
}

func TestParseGasGuide(t *testing.T) {
	path := gasListPath(t)
	entries := parseGasGuide()
	if entries == nil || len(entries) == 0 {
		t.Fatalf("expected entries from gas guide, got %v", entries)
	}
	if cost, ok := entries["Add"]; !ok || cost != 1 {
		t.Fatalf("expected Add cost 1, got %d (present=%v)", cost, ok)
	}

	// Rename the guide to simulate missing file behaviour.
	backup := path + ".bak"
	if err := os.Rename(path, backup); err != nil {
		t.Fatalf("rename gas guide: %v", err)
	}
	defer os.Rename(backup, path)
	if m := parseGasGuide(); m != nil {
		t.Fatalf("expected nil map when guide missing, got %v", m)
	}
}

func TestDefaultGasTableOverrides(t *testing.T) {
	path := gasListPath(t)
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read gas guide: %v", err)
	}
	t.Cleanup(func() {
		// Restore original file and gas table state.
		if err := os.WriteFile(path, original, 0644); err != nil {
			t.Fatalf("restore gas guide: %v", err)
		}
		initGasTable()
	})

	override := []byte("| Function | Gas Cost |\n|---|---|\n| `Add` | `5` |\n")
	if err := os.WriteFile(path, override, 0644); err != nil {
		t.Fatalf("write override gas guide: %v", err)
	}

	tbl := DefaultGasTable()
	addOp, ok := nameToOp["Add"]
	if !ok {
		t.Fatalf("Add opcode not found in catalogue")
	}
	if tbl[addOp] != 5 {
		t.Fatalf("expected override cost 5 for Add, got %d", tbl[addOp])
	}

	subOp, ok := nameToOp["Sub"]
	if !ok {
		t.Fatalf("Sub opcode not found in catalogue")
	}
	if tbl[subOp] != DefaultGasCost {
		t.Fatalf("expected default cost %d for Sub, got %d", DefaultGasCost, tbl[subOp])
	}
}

func TestSetGasCostAndSnapshot(t *testing.T) {
	initGasTable()
	op, ok := nameToOp["Add"]
	if !ok {
		t.Fatalf("Add opcode not found")
	}
	original := GasCost(op)
	SetGasCost(op, original+10)
	if GasCost(op) != original+10 {
		t.Fatalf("expected gas cost %d, got %d", original+10, GasCost(op))
	}

	snap := GasTableSnapshot()
	if snap[op] != original+10 {
		t.Fatalf("snapshot has cost %d, want %d", snap[op], original+10)
	}
	snap[op] = 0
	if GasCost(op) != original+10 {
		t.Fatalf("modifying snapshot should not alter gas table")
	}
}

func TestAccessControlGasCosts(t *testing.T) {
	initGasTable()
	grantOp, ok := nameToOp["GrantRole"]
	if !ok {
		t.Fatalf("GrantRole opcode missing")
	}
	if GasCost(grantOp) != 100 {
		t.Fatalf("expected GrantRole cost 100, got %d", GasCost(grantOp))
	}
	hasOp, ok := nameToOp["HasRole"]
	if !ok {
		t.Fatalf("HasRole opcode missing")
	}
	if GasCost(hasOp) != 30 {
		t.Fatalf("expected HasRole cost 30, got %d", GasCost(hasOp))
	}
}

func TestGasCostByName(t *testing.T) {
	initGasTable()
	if c := GasCostByName("Add"); c == 0 {
		t.Fatalf("expected non-zero cost for Add")
	}
	if c := GasCostByName("NotARealOp"); c != DefaultGasCost {
		t.Fatalf("expected default cost %d for unknown, got %d", DefaultGasCost, c)
	}
}
