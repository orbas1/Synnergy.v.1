package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"synnergy/core"
)

func execCommand(args ...string) (string, error) {
	cmd := RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	cmd.SetArgs([]string{})
	w.Close()
	os.Stdout = old
	outBuf, _ := io.ReadAll(r)
	r.Close()
	return strings.TrimSpace(buf.String() + string(outBuf)), err
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

// TestRootHelp ensures the root command prints usage information.
func TestRootHelp(t *testing.T) {
	out, err := execCommand("--help")
	if err != nil {
		t.Fatalf("help failed: %v", err)
	}
	if !strings.Contains(out, "Synnergy blockchain CLI") {
		t.Fatalf("unexpected help output: %s", out)
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

func TestConsensusMineGas(t *testing.T) {
	out, err := execCommand("consensus", "mine", "1")
	if err != nil {
		t.Fatalf("mine failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") {
		t.Fatalf("expected gas cost, got %q", out)
	}
}

func TestDAOCreationGas(t *testing.T) {
	out, err := execCommand("dao", "create", "testdao", "creator")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") {
		t.Fatalf("expected gas cost, got %q", out)
	}
}
