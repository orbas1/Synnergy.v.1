package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateTables(t *testing.T) {
	dir := t.TempDir()
	funcs := filepath.Join(dir, "functions_list.txt")
	if err := os.WriteFile(funcs, []byte("pkg/file.go: func Example()\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	op := filepath.Join(dir, "op.md")
	gas := filepath.Join(dir, "gas.md")
	if err := generateTables(funcs, op, gas, 0x100000); err != nil {
		t.Fatalf("generateTables: %v", err)
	}
	if _, err := os.Stat(op); err != nil {
		t.Fatalf("opcode file not created: %v", err)
	}
	if _, err := os.Stat(gas); err != nil {
		t.Fatalf("gas file not created: %v", err)
	}
}
