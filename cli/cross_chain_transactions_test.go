package cli

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"synnergy/core"
)

func TestCrossChainTransactionsCommands(t *testing.T) {
	l := core.NewLedger()
	l.Credit("alice", 100)
	l.Credit("bob", 100)
	crossTxManager = core.NewCrossChainTxManager(l)

	// ensure json flag defaults to false
	if cmd, _, err := rootCmd.Find([]string{"cross_tx"}); err == nil {
		_ = cmd.PersistentFlags().Set("json", "false")
	}

	out, err := execCLI(t, "cross_tx", "lockmint", "1", "asset1", "40", "proof", "--from", "alice", "--to", "carol")
	if err != nil {
		t.Fatalf("lockmint: %v", err)
	}
	fields := strings.Fields(out)
	if len(fields) < 1 {
		t.Fatalf("unexpected output: %q", out)
	}
	id1, err := strconv.Atoi(fields[0])
	if err != nil {
		t.Fatalf("invalid id: %v", err)
	}

	out, err = execCLI(t, "cross_tx", "burnrelease", "1", "dave", "asset1", "30", "--from", "bob")
	if err != nil {
		t.Fatalf("burnrelease: %v", err)
	}
	fields = strings.Fields(out)
	if len(fields) < 1 {
		t.Fatalf("unexpected output: %q", out)
	}
	if _, err := strconv.Atoi(fields[0]); err != nil {
		t.Fatalf("invalid id: %v", err)
	}

	// list transfers
	out, err = execCLI(t, "cross_tx", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "bridge=1") {
		t.Fatalf("expected bridge records, got %q", out)
	}

	// get as JSON
	out, err = execCLI(t, "cross_tx", "--json", "get", strconv.Itoa(id1))
	if err != nil {
		t.Fatalf("get json: %v", err)
	}
	var tr core.CrossChainTransfer
	if err := json.Unmarshal([]byte(out), &tr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if tr.Amount != 40 || tr.From != "alice" || tr.To != "carol" {
		t.Fatalf("unexpected transfer %+v", tr)
	}

	if bal := l.GetBalance("alice"); bal != 60 {
		t.Fatalf("expected alice balance 60, got %d", bal)
	}
	if bal := l.GetBalance("carol"); bal != 40 {
		t.Fatalf("expected carol balance 40, got %d", bal)
	}
}
