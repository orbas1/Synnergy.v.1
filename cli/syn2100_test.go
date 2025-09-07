package cli

import (
	"testing"

	"synnergy/core"
)

func run2100(t *testing.T, args ...string) {
	t.Helper()
	cmd := RootCmd()
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("command %v failed: %v", args, err)
	}
}

func TestSyn2100RegisterAndFinance(t *testing.T) {
	syn2100 = core.NewTradeFinanceToken()
	run2100(t, "syn2100", "register", "--id", "doc1", "--issuer", "A", "--recipient", "B", "--amount", "100")
	if _, ok := syn2100.GetDocument("doc1"); !ok {
		t.Fatalf("document not registered")
	}
	run2100(t, "syn2100", "finance", "doc1", "fin1")
	d, ok := syn2100.GetDocument("doc1")
	if !ok || !d.Financed || d.Financier != "fin1" {
		t.Fatalf("document not financed correctly")
	}
	run2100(t, "syn2100", "add-liquidity", "addr1", "50")
	if syn2100.Liquidity["addr1"] != 50 {
		t.Fatalf("liquidity not added")
	}
	run2100(t, "syn2100", "remove-liquidity", "addr1", "20")
	if syn2100.Liquidity["addr1"] != 30 {
		t.Fatalf("liquidity not removed correctly")
	}
}
