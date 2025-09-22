package cli

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"synnergy/core"
)

type memoryWallet struct {
	wallet   *core.Wallet
	password string
}

var (
	memoryWalletsMu sync.Mutex
	memoryWallets   map[string]memoryWallet
	memorySeq       uint64
	stage73Path     string
	stage73Once     sync.Once
)

// setStage73StatePath allows stage 74 CLI tests to override the path used by
// stage 73 stateful helpers that live in other packages. Earlier stages stored
// the value in a separate file which is not built in this repository, so the
// tests in this workspace stub the behaviour.
func setStage73StatePath(path string) {
	memoryWalletsMu.Lock()
	stage73Path = path
	memoryWalletsMu.Unlock()
}

// resetStage73LoadedForTests clears the once guard so tests can reinitialise
// stage 73 fixtures after changing the path. The real implementation lived in
// stage 73 sources, so we replicate the minimal behaviour to keep later stages
// compiling.
func resetStage73LoadedForTests() {
	memoryWalletsMu.Lock()
	stage73Once = sync.Once{}
	memoryWalletsMu.Unlock()
}

func useMemoryWalletLoader(t *testing.T) {
	t.Helper()
	memoryWalletsMu.Lock()
	if memoryWallets == nil {
		memoryWallets = make(map[string]memoryWallet)
	} else {
		for k := range memoryWallets {
			delete(memoryWallets, k)
		}
	}
	memorySeq = 0
	prev := loadWallet
	loadWallet = func(path, password string) (*core.Wallet, error) {
		if !strings.HasPrefix(path, "memory-wallet-") {
			return prev(path, password)
		}
		memoryWalletsMu.Lock()
		defer memoryWalletsMu.Unlock()
		entry, ok := memoryWallets[path]
		if !ok {
			return nil, fmt.Errorf("wallet not found: %s", path)
		}
		if entry.password != password {
			return nil, fmt.Errorf("cipher: message authentication failed")
		}
		return entry.wallet, nil
	}
	memoryWalletsMu.Unlock()
	t.Cleanup(func() {
		memoryWalletsMu.Lock()
		loadWallet = prev
		for k := range memoryWallets {
			delete(memoryWallets, k)
		}
		memoryWalletsMu.Unlock()
	})
}

func newMemoryWallet(t *testing.T, password string) (*core.Wallet, string) {
	t.Helper()
	w, err := core.NewWallet()
	if err != nil {
		t.Fatalf("wallet: %v", err)
	}
	memoryWalletsMu.Lock()
	defer memoryWalletsMu.Unlock()
	id := fmt.Sprintf("memory-wallet-%d", memorySeq)
	memorySeq++
	memoryWallets[id] = memoryWallet{wallet: w, password: password}
	return w, id
}
