package cli

import (
	"testing"

	"synnergy/core"
)

func run2500(args ...string) error {
	cmd := RootCmd()
	cmd.SetArgs(args)
	return cmd.Execute()
}

func TestSyn2500AddAndUpdate(t *testing.T) {
	syn2500 = core.NewSyn2500Registry()
	if err := run2500("syn2500", "add", "--id", "m1", "--addr", "a", "--power", "10"); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	m, ok := syn2500.GetMember("m1")
	if !ok || m.VotingPower != 10 {
		t.Fatalf("member not added correctly")
	}
	if err := run2500("syn2500", "update", "m1", "20"); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	m, _ = syn2500.GetMember("m1")
	if m.VotingPower != 20 {
		t.Fatalf("voting power not updated")
	}
}

func TestSyn2500AddMissingFields(t *testing.T) {
	syn2500 = core.NewSyn2500Registry()
	if err := run2500("syn2500", "add", "--id", "", "--addr", "", "--power", "0"); err == nil {
		t.Fatalf("expected error for missing fields")
	}
}
