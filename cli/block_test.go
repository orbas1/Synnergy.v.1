package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

// TestBlockCreateAndHeader ensures block creation and header hashing work via CLI.
func TestBlockCreateAndHeader(t *testing.T) {
	sbList = nil
	lastBlock = nil

	wallet, err := core.NewWallet()
	if err != nil {
		t.Fatalf("wallet: %v", err)
	}
	if err := core.RegisterValidatorWallet(wallet); err != nil {
		t.Fatalf("register: %v", err)
	}
	t.Cleanup(func() { core.UnregisterValidator(wallet.Address) })

	if _, err := execCommand("block", "sub-create", wallet.Address, "from", "to", "10", "1", "0"); err != nil {
		t.Fatalf("sub-create failed: %v", err)
	}

	out, err := execCommand("block", "create", "prevhash")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if !strings.Contains(out, "block with 1 sub-blocks") {
		t.Fatalf("unexpected create output: %s", out)
	}

	out, err = execCommand("block", "header", "1")
	if err != nil {
		t.Fatalf("header failed: %v", err)
	}
	if out == "" {
		t.Fatalf("expected header hash")
	}
}
