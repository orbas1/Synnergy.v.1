package core

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
	mu           sync.RWMutex
	contracts    map[string]*Contract
	vm           VirtualMachine
	ledger       *Ledger
	feeCollector string
}

func contractFeeCollectorAddress() string {
	hash := sha256.Sum256([]byte("contract_fee_pool"))
	return hex.EncodeToString(hash[:])
}

// NewContractRegistry initialises an empty registry backed by the provided VM.
// If a ledger is provided any persisted contracts are preloaded into memory.
func NewContractRegistry(vm VirtualMachine, ledger *Ledger) *ContractRegistry {
	reg := &ContractRegistry{
		contracts:    make(map[string]*Contract),
		vm:           vm,
		ledger:       ledger,
		feeCollector: contractFeeCollectorAddress(),
	}
	if ledger != nil {
		for _, rec := range ledger.Contracts() {
			wasm := make([]byte, len(rec.WASM))
			copy(wasm, rec.WASM)
			reg.contracts[rec.Address] = &Contract{
				Address:  rec.Address,
				Owner:    rec.Owner,
				WASM:     wasm,
				Manifest: rec.Manifest,
				GasLimit: rec.GasLimit,
			}
		}
	}
	return reg
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

	r.mu.RLock()
	_, exists := r.contracts[addr]
	r.mu.RUnlock()
	if exists {
		return "", errors.New("contract already deployed")
	}
	if r.ledger != nil && gasLimit > 0 {
		if err := r.ledger.Transfer(owner, r.feeCollector, gasLimit, 0); err != nil {
			return "", err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.contracts[addr]; exists {
		if r.ledger != nil && gasLimit > 0 {
			_ = r.ledger.Transfer(r.feeCollector, owner, gasLimit, 0)
		}
		return "", errors.New("contract already deployed")
	}
	r.contracts[addr] = &Contract{
		Address:  addr,
		Owner:    owner,
		WASM:     wasm,
		Manifest: manifest,
		GasLimit: gasLimit,
	}
	if r.ledger != nil {
		stored := make([]byte, len(wasm))
		copy(stored, wasm)
		r.ledger.RegisterContract(LedgerContract{
			Address:  addr,
			Owner:    owner,
			Manifest: manifest,
			GasLimit: gasLimit,
			WASM:     stored,
		})
	}
	return addr, nil
}

// Invoke executes a method on the specified contract via the configured VM.
// It returns the output bytes and the gas consumed.
func (r *ContractRegistry) Invoke(addr, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {
	return r.InvokeFrom(addr, "", method, args, gasLimit)
}

// InvokeFrom executes a method on the specified contract, charging gas to the
// supplied caller. When caller is empty the contract owner is billed.
func (r *ContractRegistry) InvokeFrom(addr, caller, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {
	r.mu.RLock()
	c, ok := r.contracts[addr]
	r.mu.RUnlock()
	if !ok {
		return nil, 0, errors.New("contract not found")
	}
	if c.Paused {
		return nil, 0, errors.New("contract paused")
	}
	limit := gasLimit
	if limit == 0 || limit > c.GasLimit {
		limit = c.GasLimit
	}
	payer := caller
	if payer == "" {
		payer = c.Owner
	}
	if r.ledger != nil && limit > 0 {
		if err := r.ledger.Transfer(payer, r.feeCollector, limit, 0); err != nil {
			return nil, 0, err
		}
	}
	out, used, err := r.vm.Execute(c.WASM, method, args, limit)
	if err != nil {
		if r.ledger != nil && limit > 0 {
			_ = r.ledger.Transfer(r.feeCollector, payer, limit, 0)
		}
		return out, used, err
	}
	if r.ledger != nil && used < limit {
		refund := limit - used
		if refund > 0 {
			_ = r.ledger.Transfer(r.feeCollector, payer, refund, 0)
		}
	}
	return out, used, nil
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
