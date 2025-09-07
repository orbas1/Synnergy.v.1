package cli

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"synnergy/core"
)

func resetSyn4900() { agriRegistry = core.NewAgriculturalRegistry() }

func TestTokenSyn4900RegisterInfo(t *testing.T) {
	resetSyn4900()
	harvest := strconv.FormatInt(time.Now().Unix(), 10)
	expiry := strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10)
	if _, err := execCommand("token_syn4900", "register", "A1", "grain", "alice", "US", "100", harvest, expiry, "cert"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	out, err := execCommand("token_syn4900", "--json", "info", "A1")
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "A1") {
		t.Fatalf("unexpected output: %s", out)
	}
}
