package synnergy

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
)

// Contract represents a deployed smart contract. It keeps minimal metadata
// required by the CLI and other modules to manage and invoke contracts.
type Contract struct {
	Address  string // deterministic hash of the WASM bytecode
	Owner    string // creator/owner address
	WASM     []byte // raw WASM bytecode
	Manifest string // optional Ricardian manifest JSON
	GasLimit uint64 // max gas allowed per invocation
	Paused   bool   // whether execution is paused
}

// VirtualMachine defines the execution interface required by the contract
// registry. The VM is implemented in virtual_machine.go.
type VirtualMachine interface {
	Execute(wasm []byte, method string, args []byte, gasLimit uint64) ([]byte, uint64, error)
	Start() error
	Stop() error
	Status() bool
}

// ContractRegistry stores deployed contracts and offers helper methods for
// deployment and invocation. It is safe for concurrent use.
type ContractRegistry struct {
	mu        sync.RWMutex
	contracts map[string]*Contract
	vm        VirtualMachine
}

// NewContractRegistry initialises an empty registry backed by the provided VM.
func NewContractRegistry(vm VirtualMachine) *ContractRegistry {
	return &ContractRegistry{
		contracts: make(map[string]*Contract),
		vm:        vm,
	}
}

// CompileWASM returns the input bytecode and its sha256 hash. In the full
// implementation this would also convert WAT to WASM, but here we simply return
// the bytes unmodified for deterministic builds.
func CompileWASM(src []byte) ([]byte, string, error) {
	if len(src) == 0 {
		return nil, "", errors.New("source bytecode is empty")
	}
	h := sha256.Sum256(src)
	return src, hex.EncodeToString(h[:]), nil
}

// Deploy registers a new contract. The address is derived from the bytecode
// hash. If a contract with the same address already exists an error is
// returned.
func (r *ContractRegistry) Deploy(wasm []byte, manifest string, gasLimit uint64, owner string) (string, error) {
	if len(wasm) == 0 {
		return "", errors.New("wasm bytecode required")
	}
	hash := sha256.Sum256(wasm)
	addr := hex.EncodeToString(hash[:])

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.contracts[addr]; exists {
		return "", errors.New("contract already deployed")
	}
	r.contracts[addr] = &Contract{
		Address:  addr,
		Owner:    owner,
		WASM:     wasm,
		Manifest: manifest,
		GasLimit: gasLimit,
	}
	return addr, nil
}

// Invoke executes a method on the specified contract via the configured VM.
// It returns the output bytes and the gas consumed.
func (r *ContractRegistry) Invoke(addr, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {
	r.mu.RLock()
	c, ok := r.contracts[addr]
	r.mu.RUnlock()
	if !ok {
		return nil, 0, errors.New("contract not found")
	}
	if c.Paused {
		return nil, 0, errors.New("contract paused")
	}
	if gasLimit == 0 || gasLimit > c.GasLimit {
		gasLimit = c.GasLimit
	}
	return r.vm.Execute(c.WASM, method, args, gasLimit)
}

// List returns all deployed contracts.
func (r *ContractRegistry) List() []*Contract {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Contract, 0, len(r.contracts))
	for _, c := range r.contracts {
		out = append(out, c)
	}
	return out
}

// Get fetches a contract by address.
func (r *ContractRegistry) Get(addr string) (*Contract, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.contracts[addr]
	return c, ok
}
