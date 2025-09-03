package cli

import (
	"bytes"
	"strings"
	"testing"

	"synnergy/core"
)

func execCommand(args ...string) (string, error) {
	cmd := RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	cmd.SetArgs([]string{})
	return strings.TrimSpace(buf.String()), err
}

func TestAddressParse(t *testing.T) {
	out, err := execCommand("address", "parse", core.AddressZero.Hex())
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if out != core.AddressZero.Hex() {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestNetworkStartStop(t *testing.T) {
	network.Stop()
	if _, err := execCommand("network", "start"); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if _, err := execCommand("network", "stop"); err != nil {
		t.Fatalf("stop failed: %v", err)
	}
}

func TestPeerDiscoverEmpty(t *testing.T) {
	network.Start()
	defer network.Stop()
	out, err := execCommand("peer", "discover")
	if err != nil {
		t.Fatalf("discover failed: %v", err)
	}
	if out != "" {
		t.Fatalf("expected no peers, got %q", out)
	}
}
