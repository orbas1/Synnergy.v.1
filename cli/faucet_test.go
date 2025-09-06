package cli

import (
	"strings"
	"testing"
)

func TestFaucetRequest(t *testing.T) {
	if _, err := execCommand("faucet", "init", "--balance", "2", "--amount", "1", "--cooldown", "0s"); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	out, err := execCommand("faucet", "request", "addr1")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if !strings.Contains(out, "dispensed") {
		t.Fatalf("unexpected output: %s", out)
	}
}
