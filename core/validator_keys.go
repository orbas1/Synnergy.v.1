package core

import (
	"crypto/ecdsa"
	"errors"
	"sync"
)

// validatorKeyStore keeps track of validator signing keys used for sub-block
// signatures. Keys are registered by nodes when validators join the network so
// that all components can authenticate sub-blocks without embedding private
// key material in the structures themselves.
type validatorKeyStore struct {
	mu       sync.RWMutex
	privKeys map[string]*ecdsa.PrivateKey
	pubKeys  map[string]*ecdsa.PublicKey
}

var globalValidatorKeys = newValidatorKeyStore()

func newValidatorKeyStore() *validatorKeyStore {
	return &validatorKeyStore{
		privKeys: make(map[string]*ecdsa.PrivateKey),
		pubKeys:  make(map[string]*ecdsa.PublicKey),
	}
}

// RegisterValidatorWallet records the wallet keys for the given validator
// address. The wallet's public key is shared with other components, while the
// private key is retained only to support local signing helpers used by tests
// and CLI tooling.
func RegisterValidatorWallet(w *Wallet) error {
	if w == nil || w.PrivateKey == nil {
		return errors.New("wallet private key not initialised")
	}
	addr := w.Address
	if addr == "" {
		return errors.New("wallet address required")
	}
	globalValidatorKeys.mu.Lock()
	globalValidatorKeys.privKeys[addr] = w.PrivateKey
	pub := w.PrivateKey.PublicKey
	globalValidatorKeys.pubKeys[addr] = &pub
	globalValidatorKeys.mu.Unlock()
	return nil
}

// UnregisterValidator removes any stored keys for the provided address. It is
// primarily used by tests to ensure isolation between cases.
func UnregisterValidator(addr string) {
	globalValidatorKeys.mu.Lock()
	delete(globalValidatorKeys.privKeys, addr)
	delete(globalValidatorKeys.pubKeys, addr)
	globalValidatorKeys.mu.Unlock()
}

func validatorPrivateKey(addr string) *ecdsa.PrivateKey {
	globalValidatorKeys.mu.RLock()
	defer globalValidatorKeys.mu.RUnlock()
	return globalValidatorKeys.privKeys[addr]
}

func validatorPublicKey(addr string) *ecdsa.PublicKey {
	globalValidatorKeys.mu.RLock()
	defer globalValidatorKeys.mu.RUnlock()
	if pub, ok := globalValidatorKeys.pubKeys[addr]; ok {
		return pub
	}
	return nil
}
