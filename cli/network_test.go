package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"synnergy/core"
)

// TestNetworkCLIStartStop verifies start and stop commands emit output.
func TestNetworkCLIStartStop(t *testing.T) {
	network = core.NewNetwork(biometricSvc)
	network.Stop() // ensure known state

	out, err := execNetCLI("network", "start")
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if !strings.Contains(out, "network started") {
		t.Fatalf("unexpected output: %s", out)
	}

	out, err = execNetCLI("network", "stop")
	if err != nil {
		t.Fatalf("stop failed: %v", err)
	}
	if !strings.Contains(out, "network stopped") {
		t.Fatalf("unexpected output: %s", out)
	}
}

// TestNetworkCLIPeersAndBroadcast covers peer listing and broadcasting.
func TestNetworkCLIPeersAndBroadcast(t *testing.T) {
	network = core.NewNetwork(biometricSvc)
	network.AddNode(core.NewNode("n1", "addr", core.NewLedger()))

	out, err := execNetCLI("--json", "network", "peers")
	if err != nil {
		t.Fatalf("peers failed: %v", err)
	}
	if !strings.Contains(out, "n1") {
		t.Fatalf("expected peer in output: %s", out)
	}

	ch := network.Subscribe("topic")
	if _, err := execNetCLI("network", "broadcast", "topic", "hi"); err != nil {
		t.Fatalf("broadcast failed: %v", err)
	}
	select {
	case msg := <-ch:
		if string(msg) != "hi" {
			t.Fatalf("unexpected message: %s", string(msg))
		}
	case <-time.After(time.Second):
		t.Fatal("no message received")
	}

	network.Stop()
}

// execNetCLI executes rootCmd with args while capturing stdout.
func execNetCLI(args ...string) (string, error) {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetErr(new(bytes.Buffer))
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	w.Close()
	os.Stdout = stdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}
