package cli

import "testing"

func run300(args ...string) error {
	cmd := RootCmd()
	cmd.SetArgs(args)
	return cmd.Execute()
}

func TestSyn300InitAndProposal(t *testing.T) {
	syn300 = nil
	if err := run300("syn300", "init", "--balances", "alice=100"); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	if err := run300("syn300", "propose", "alice", "test"); err != nil {
		t.Fatalf("propose failed: %v", err)
	}
	if err := run300("syn300", "vote", "1", "alice", "true"); err != nil {
		t.Fatalf("vote failed: %v", err)
	}
	if err := run300("syn300", "execute", "1", "100"); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
}

func TestSyn300InitMissingBalances(t *testing.T) {
	syn300 = nil
	if err := run300("syn300", "init", "--balances", ""); err == nil {
		t.Fatalf("expected error for missing balances")
	}
}
