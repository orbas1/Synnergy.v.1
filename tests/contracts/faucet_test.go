package contracts

import (
        "bytes"
        "io"
        "os"
        "strings"
        "testing"

        synn "synnergy"
        "synnergy/cli"
)

// execCLI executes the Synnergy CLI command with the given arguments from the
// repository root and returns combined stdout and stderr.
func execCLI(t *testing.T, args ...string) (string, error) {
        t.Helper()
        wd, err := os.Getwd()
        if err != nil {
                t.Fatalf("getwd: %v", err)
        }
        defer os.Chdir(wd)
        if err := os.Chdir("../.." ); err != nil {
                t.Fatalf("chdir: %v", err)
        }
        cmd := cli.RootCmd()
        buf := new(bytes.Buffer)
        cmd.SetOut(buf)
        cmd.SetErr(buf)
        r, w, _ := os.Pipe()
        old := os.Stdout
        os.Stdout = w
        cmd.SetArgs(args)
        _, err = cmd.ExecuteC()
        cmd.SetArgs([]string{})
        w.Close()
        os.Stdout = old
        out, _ := io.ReadAll(r)
        r.Close()
        return strings.TrimSpace(buf.String() + string(out)), err
}

// TestTokenFaucetTemplate deploys the token faucet contract template via the
// CLI and ensures it is registered.
func TestTokenFaucetTemplate(t *testing.T) {
        synn.LoadGasTable()
        out, err := execCLI(t, "contracts", "deploy-template", "--name", "token_faucet")
        if err != nil {
                t.Fatalf("deploy-template: %v", err)
        }
        addr := strings.TrimSpace(out)
        if addr == "" {
                t.Fatal("empty address returned")
        }
        list, err := execCLI(t, "contracts", "list")
        if err != nil {
                t.Fatalf("list: %v", err)
        }
        if !strings.Contains(list, addr) {
                t.Fatalf("deployed address %s not listed", addr)
        }
}

