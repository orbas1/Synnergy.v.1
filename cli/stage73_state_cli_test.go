package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStage73StateCorruptFile(t *testing.T) {
	useMemoryWalletLoader(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "stage73.json")
	if err := os.WriteFile(path, []byte("{\"index\":"), 0o600); err != nil {
		t.Fatalf("write corrupt snapshot: %v", err)
	}
	setStage73StatePath(path)
	resetStage73LoadedForTests()
	t.Cleanup(func() {
		setStage73StatePath("")
		resetStage73LoadedForTests()
	})
	if _, err := execCommand("syn3700", "snapshot"); err == nil {
		t.Fatal("expected snapshot load failure")
	} else if !strings.Contains(err.Error(), "decode snapshot") {
		t.Fatalf("unexpected error: %v", err)
	}
}
