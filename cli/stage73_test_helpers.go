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
)

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
