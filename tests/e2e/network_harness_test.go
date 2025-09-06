package tests

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
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

// execCLI executes the Synnergy CLI command and returns combined output.
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
	outStr := strings.TrimSpace(buf.String() + string(out))
	if start := strings.IndexAny(outStr, "{"); start != -1 {
		if end := strings.LastIndex(outStr, "}"); end != -1 && end >= start {
			outStr = outStr[start : end+1]
		}
	} else if start := strings.Index(outStr, "["); start != -1 {
		if end := strings.LastIndex(outStr, "]"); end != -1 && end >= start {
			outStr = outStr[start : end+1]
		}
	}
	return outStr, err
}

// TestNetworkHarness spins up core components and exercises a full flow
// from wallet creation through transaction broadcast, demonstrating
// interoperability between CLI, wallet server and in-memory nodes.
func TestNetworkHarness(t *testing.T) {
	synn.LoadGasTable()

	// Start wallet server in a subprocess so HTTP endpoints are available
	srv := exec.Command("go", "run", "./walletserver")
	srv.Dir = ".."
	if err := srv.Start(); err != nil {
		t.Skipf("start wallet server: %v", err)
	}
	defer srv.Process.Kill()

	// Wait for wallet server health endpoint
	for i := 0; i < 20; i++ {
		time.Sleep(200 * time.Millisecond)
		resp, err := http.Get("http://localhost:8080/health")
		if err == nil {
			resp.Body.Close()
			break
		}
		if i == 19 {
			t.Skipf("wallet server not responding: %v", err)
		}
	}

	// Create a wallet via HTTP as GUI would
	resp, err := http.Post("http://localhost:8080/wallet/new", "application/json", nil)
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

	// Use CLI to validate the wallet address
	if out, err := execCLI(t, "address", "parse", data.Address); err != nil || strings.TrimPrefix(out, "0x") != strings.TrimPrefix(data.Address, "0x") {
		t.Fatalf("cli address parse failed: %v %s", err, out)
	}

	// Build an in-memory network with biometric authentication
	svc := core.NewBiometricService()
	network := core.NewNetwork(svc)
	n1 := core.NewNode("n1", "addr1", core.NewLedger())
	n2 := core.NewNode("n2", "addr2", core.NewLedger())
	network.AddNode(n1)
	network.AddNode(n2)
	n1.Ledger.Credit(data.Address, 100)
	n2.Ledger.Credit(data.Address, 100)

	// Enrol biometric data for the new wallet
	bio := []byte("finger")
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	svc.Enroll(data.Address, bio, &key.PublicKey)
	h := sha256.Sum256(bio)
	sig, err := ecdsa.SignASN1(rand.Reader, key, h[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	// Broadcast a simple transaction through the network
	tx := core.NewTransaction(data.Address, "recipient", 1, 1, 1)
	if err := network.Broadcast(tx, data.Address, bio, sig); err != nil {
		t.Fatalf("broadcast failed: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	if len(n1.Mempool) != 1 || len(n2.Mempool) != 1 {
		t.Fatalf("transaction not received by all nodes")
	}

	// Ensure gas snapshot CLI command returns data
	out, err := execCLI(t, "gas", "snapshot", "--json")
	if err != nil {
		t.Fatalf("gas snapshot failed: %v", err)
	}
	var m map[string]uint64
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(m) == 0 {
		t.Fatalf("expected non-empty gas snapshot")
	}
}
