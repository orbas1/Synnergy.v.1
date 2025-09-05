package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestAuditLogAndList(t *testing.T) {
	addr := "0x1234567890abcdef1234567890abcdef12345678"
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"audit", "log", addr, "login"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("log command failed: %v", err)
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"audit", "list", addr})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list command failed: %v", err)
	}
	if !strings.Contains(buf.String(), "login") {
		t.Fatalf("expected login event, got %s", buf.String())
	}
}
