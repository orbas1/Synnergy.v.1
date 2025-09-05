package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestAuditNodeStart(t *testing.T) {
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"audit_node", "start"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("start command failed: %v", err)
	}
}

func TestAuditNodeLogAndList(t *testing.T) {
	addr := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"audit_node", "log", addr, "update"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("log command failed: %v", err)
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"audit_node", "list", addr})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list command failed: %v", err)
	}
	if !strings.Contains(buf.String(), "update") {
		t.Fatalf("expected update event, got %s", buf.String())
	}
}
