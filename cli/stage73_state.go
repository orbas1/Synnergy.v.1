package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	stage73Mu     sync.Mutex
	stage73Path   string
	stage73Loaded bool
)

// setStage73StatePath configures the shared Stage 73 persistence location used by
// CLI tests and the function web. Tests override the path to ensure isolation.
func setStage73StatePath(path string) {
	stage73Mu.Lock()
	defer stage73Mu.Unlock()
	stage73Path = path
	stage73Loaded = false
	if stage73Path != "" {
		_ = os.MkdirAll(filepath.Dir(stage73Path), 0o755)
	}
}

// resetStage73LoadedForTests clears the cached Stage 73 load flag so fresh state
// is materialised for each test. Any persisted file at the configured path is
// removed to avoid bleeding state across suites.
func resetStage73LoadedForTests() {
	stage73Mu.Lock()
	defer stage73Mu.Unlock()
	stage73Loaded = false
	if stage73Path != "" {
		_ = os.Remove(stage73Path)
	}
}

func stage73StatePath() string {
	stage73Mu.Lock()
	defer stage73Mu.Unlock()
	return stage73Path
}

func markStage73Loaded() {
	stage73Mu.Lock()
	stage73Loaded = true
	stage73Mu.Unlock()
}

func stage73HasLoaded() bool {
	stage73Mu.Lock()
	defer stage73Mu.Unlock()
	return stage73Loaded
}

// parseWalletCredential expects values in the form "path:password" so commands
// can load wallets without interactive prompts.
func parseWalletCredential(value string) (string, string, error) {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid wallet credential format")
	}
	path := strings.TrimSpace(parts[0])
	password := parts[1]
	if path == "" {
		return "", "", fmt.Errorf("wallet path required")
	}
	if password == "" {
		return "", "", fmt.Errorf("wallet password required")
	}
	return path, password, nil
}
