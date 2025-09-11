package cli

import (
	"os"
	"strings"
	"testing"
)

func TestGuiCommandRegistered(t *testing.T) {
	cmd := RootCmd()
	// Ensure the gui command exists
	found := false
	for _, c := range cmd.Commands() {
		if c.Name() == "gui" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("gui command not registered")
	}

	// Ensure help executes without launching the GUI
	cmd.SetArgs([]string{"gui", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error executing help: %v", err)
	}
}

func TestGuiCommandMissingDirectory(t *testing.T) {
	tmp := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(cwd)
	os.Chdir(tmp)
	if err := guiCmd.RunE(guiCmd, nil); err == nil {
		t.Fatal("expected error when GUI directory missing")
	} else if !strings.Contains(err.Error(), "desktop GUI shell not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}
