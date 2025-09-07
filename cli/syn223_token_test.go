package cli

import (
	"testing"
)

func run223(t *testing.T, args ...string) {
	t.Helper()
	cmd := RootCmd()
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("command %v failed: %v", args, err)
	}
}

func TestSyn223InitTransfer(t *testing.T) {
	syn223 = nil
	run223(t, "syn223", "init", "--name", "Token", "--symbol", "T", "--owner", "alice", "--supply", "1000")
	if syn223 == nil || syn223.BalanceOf("alice") != 1000 {
		t.Fatalf("initialisation failed")
	}
	run223(t, "syn223", "whitelist", "bob")
	run223(t, "syn223", "transfer", "alice", "bob", "200")
	if syn223.BalanceOf("bob") != 200 {
		t.Fatalf("transfer failed")
	}
	run223(t, "syn223", "balance", "bob")
	if syn223.BalanceOf("bob") != 200 {
		t.Fatalf("balance command failed")
	}
}
