package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	pipeBuf := new(bytes.Buffer)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, _ = io.Copy(pipeBuf, r)
		wg.Done()
	}()
	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	cmd.SetArgs([]string{})
	resetCommandFlags(cmd)
	w.Close()
	wg.Wait()
	os.Stdout = old
	r.Close()
	return strings.TrimSpace(buf.String() + pipeBuf.String()), err
}

func resetCommandFlags(cmd *cobra.Command) {
	resetFlagSet(cmd.Flags())
	resetFlagSet(cmd.PersistentFlags())
	for _, c := range cmd.Commands() {
		resetCommandFlags(c)
	}
}

func resetFlagSet(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}
	fs.VisitAll(func(f *pflag.Flag) {
		switch f.Value.Type() {
		case "stringSlice":
			_ = f.Value.Set("")
		default:
			_ = f.Value.Set(f.DefValue)
		}
		f.Changed = false
	})
}

func jsonPayload(out string) string {
	lines := strings.Split(out, "\n")
	start := -1
	end := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if start == -1 && (strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")) {
			start = i
		}
		if strings.HasSuffix(trimmed, "}") || strings.HasSuffix(trimmed, "]") {
			end = i
		}
	}
	if start == -1 {
		return out
	}
	if end == -1 {
		end = len(lines) - 1
	}
	return strings.TrimSpace(strings.Join(lines[start:end+1], "\n"))
}

func firstNonGasLine(out string) string {
	for _, line := range strings.Split(out, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "gas cost:") {
			continue
		}
		return trimmed
	}
	return ""
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
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "[]") {
		t.Fatalf("expected gas cost and empty result, got %q", out)
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
