package cli

import (
        "bytes"
        "io"
        "os"
        "strings"
        "testing"
)

// executeCLICommand runs the root command with the provided arguments and
// captures output from both Cobra's configured writer and direct stdout
// writes used by some commands.
func executeCLICommand(t *testing.T, args ...string) (string, error) {
        t.Helper()
        cmd := RootCmd()
        buf := new(bytes.Buffer)
        cmd.SetOut(buf)
        cmd.SetErr(buf)
        cmd.SetArgs(args)
        defer cmd.SetArgs(nil)

        oldStdout := os.Stdout
        r, w, err := os.Pipe()
        if err != nil {
                t.Fatalf("pipe: %v", err)
        }
        os.Stdout = w
        execErr := cmd.Execute()
        w.Close()
        os.Stdout = oldStdout
        stdoutBytes, _ := io.ReadAll(r)
        r.Close()

        output := buf.String() + string(stdoutBytes)
        return strings.TrimSpace(output), execErr
}
