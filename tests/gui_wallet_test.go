package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	synn "synnergy"
	"synnergy/cli"
	"synnergy/core"
)

// execCLI executes the Synnergy CLI command with given args and returns output.
func execCLI(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := cli.RootCmd()
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
	out, _ := io.ReadAll(r)
	r.Close()
	return strings.TrimSpace(buf.String() + string(out)), err
}

// TestGUIWalletIntegration spins up the wallet server and exercises its HTTP
// endpoints similarly to how the GUI would, verifying the responses can be
// consumed by the CLI.
func TestGUIWalletIntegration(t *testing.T) {
	synn.LoadGasTable()

	// Start wallet server as external process.
	srv := exec.Command("go", "run", "./walletserver")
	srv.Dir = ".."
	if err := srv.Start(); err != nil {
		t.Fatalf("start wallet server: %v", err)
	}
	defer srv.Process.Kill()

	// Wait for server to become healthy.
	var resp *http.Response
	var err error
	for i := 0; i < 20; i++ {
		time.Sleep(200 * time.Millisecond)
		resp, err = http.Get("http://localhost:8080/health")
		if err == nil {
			resp.Body.Close()
			break
		}
	}
	if err != nil {
		t.Fatalf("wallet server not responding: %v", err)
	}

	// Request a new wallet via HTTP, mimicking GUI behavior.
	resp, err = http.Post("http://localhost:8080/wallet/new", "application/json", nil)
	if err != nil {
		t.Fatalf("create wallet: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	var data struct{ Address string }
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if data.Address == "" {
		t.Fatal("empty address returned")
	}

	// Ensure CLI understands the returned address.
	out, err := execCLI(t, "address", "parse", data.Address)
	if err != nil {
		t.Fatalf("cli parse: %v", err)
	}
	if strings.TrimPrefix(out, "0x") != strings.TrimPrefix(data.Address, "0x") {
		t.Fatalf("cli unexpected output: %s", out)
	}

	// Validate address format with core helper.
	if len(strings.TrimPrefix(data.Address, "0x")) != len(strings.TrimPrefix(core.AddressZero.Hex(), "0x")) {
		t.Fatalf("unexpected address length: %d", len(data.Address))
	}
}
