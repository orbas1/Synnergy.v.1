package cli

import (
	"testing"

	"synnergy/internal/tokens"
)

func run2369(args ...string) error {
	cmd := RootCmd()
	cmd.SetArgs(args)
	return cmd.Execute()
}

func TestSyn2369CreateAndTransfer(t *testing.T) {
	itemRegistry = tokens.NewItemRegistry()
	if err := run2369("syn2369", "create", "--owner", "alice", "--name", "sword"); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	it, ok := itemRegistry.GetItem("VI-1")
	if !ok || it.Owner != "alice" {
		t.Fatalf("item not created correctly")
	}
	if err := run2369("syn2369", "transfer", "VI-1", "bob"); err != nil {
		t.Fatalf("transfer failed: %v", err)
	}
	it, _ = itemRegistry.GetItem("VI-1")
	if it.Owner != "bob" {
		t.Fatalf("owner not updated")
	}
}

func TestSyn2369MissingFields(t *testing.T) {
	itemRegistry = tokens.NewItemRegistry()
	if err := run2369("syn2369", "create", "--owner", "", "--name", ""); err == nil {
		t.Fatalf("expected error for missing fields")
	}
}
