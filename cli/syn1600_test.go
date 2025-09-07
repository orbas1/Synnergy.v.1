package cli

import (
	"bytes"
	"strings"
	"testing"
)

// TestSyn1600InitRequiresFlags ensures init fails without required metadata.
func TestSyn1600InitRequiresFlags(t *testing.T) {
	musicToken = nil
	cmd := RootCmd()
	cmd.SetArgs([]string{"syn1600", "init", "--title", "Song", "--artist", "Artist"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing album flag")
	}
}

// TestSyn1600Workflow verifies init, update and info commands.
func TestSyn1600Workflow(t *testing.T) {
	musicToken = nil
	cmd := RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	cmd.SetArgs([]string{"syn1600", "init", "--title", "Song", "--artist", "Artist", "--album", "Album"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn1600", "update", "--title", "Remix"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	buf.Reset()
	cmd.SetArgs([]string{"syn1600", "info"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("info failed: %v", err)
	}
	if strings.TrimSpace(buf.String()) != "Remix by Artist on Album" {
		t.Fatalf("unexpected info output: %s", buf.String())
	}
}
