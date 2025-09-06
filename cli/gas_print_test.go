package cli

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestGasPrint(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	gasPrint("GasList")
	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "gas cost") {
		t.Fatalf("expected gas cost output, got %s", string(out))
	}
}
