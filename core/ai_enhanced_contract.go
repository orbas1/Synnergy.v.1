package core

import (
	"errors"
	"sync"
)

// AIContractRegistry wraps a base ContractRegistry to keep metadata specific to
// AI enhanced contracts such as the model hash used for inference.
type AIContractRegistry struct {
	base *ContractRegistry
	mu   sync.RWMutex
	meta map[string]string // contract address -> model hash
}

// NewAIContractRegistry creates a new registry using the provided base
// registry. The base registry handles deployment and invocation while this type
// tracks additional AI metadata.
func NewAIContractRegistry(base *ContractRegistry) *AIContractRegistry {
	return &AIContractRegistry{
		base: base,
		meta: make(map[string]string),
	}
}

// DeployAIContract deploys the WASM bytecode and records the associated model
// hash. The returned address can later be used to invoke the contract.
func (r *AIContractRegistry) DeployAIContract(wasm []byte, modelHash, manifest string, gasLimit uint64, owner string) (string, error) {
	addr, err := r.base.Deploy(wasm, manifest, gasLimit, owner)
	if err != nil {
		return "", err
	}
	r.mu.Lock()
	r.meta[addr] = modelHash
	r.mu.Unlock()
	return addr, nil
}

// InvokeAIContract invokes the "infer" method of the specified contract. The
// input payload is passed as arguments to the VM.
func (r *AIContractRegistry) InvokeAIContract(addr string, input []byte, gasLimit uint64) ([]byte, uint64, error) {
	if _, ok := r.meta[addr]; !ok {
		return nil, 0, errors.New("ai contract not found")
	}
	return r.base.Invoke(addr, "infer", input, gasLimit)
}

// ModelHash returns the stored model hash for the given contract.
func (r *AIContractRegistry) ModelHash(addr string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.meta[addr]
	return h, ok
}
