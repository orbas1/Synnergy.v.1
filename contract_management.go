package synnergy

import (
	"errors"
	"fmt"
	"strings"
)

// ContractManager provides administrative operations over deployed contracts.
type ContractManager struct {
	registry *ContractRegistry
}

var (
	ErrInvalidContractOwner = errors.New("new contract owner cannot be empty")
)

// NewContractManager wires a manager to an existing registry.
func NewContractManager(reg *ContractRegistry) *ContractManager {
	return &ContractManager{registry: reg}
}

// Transfer changes the owner of a contract.
func (m *ContractManager) Transfer(addr, newOwner string) error {
	newOwner = strings.TrimSpace(newOwner)
	if newOwner == "" {
		return ErrInvalidContractOwner
	}
	c, ok := m.registry.Get(addr)
	if !ok {
		return fmt.Errorf("%w: %s", ErrContractNotFound, addr)
	}
	m.registry.mu.Lock()
	c.Owner = newOwner
	snapshot := cloneContract(c)
	m.registry.mu.Unlock()
	m.registry.persistContract(snapshot)
	m.registry.emit(ContractRegistryEvent{
		Type:     ContractRegistryEventTransfer,
		Contract: snapshot,
		Caller:   newOwner,
	})
	return nil
}

// Pause disables contract execution.
func (m *ContractManager) Pause(addr string) error {
	c, ok := m.registry.Get(addr)
	if !ok {
		return fmt.Errorf("%w: %s", ErrContractNotFound, addr)
	}
	m.registry.mu.Lock()
	c.Paused = true
	snapshot := cloneContract(c)
	m.registry.mu.Unlock()
	m.registry.emit(ContractRegistryEvent{
		Type:     ContractRegistryEventPause,
		Contract: snapshot,
		Caller:   snapshot.Owner,
	})
	return nil
}

// Resume enables execution for a paused contract.
func (m *ContractManager) Resume(addr string) error {
	c, ok := m.registry.Get(addr)
	if !ok {
		return fmt.Errorf("%w: %s", ErrContractNotFound, addr)
	}
	m.registry.mu.Lock()
	c.Paused = false
	snapshot := cloneContract(c)
	m.registry.mu.Unlock()
	m.registry.emit(ContractRegistryEvent{
		Type:     ContractRegistryEventResume,
		Contract: snapshot,
		Caller:   snapshot.Owner,
	})
	return nil
}

// Upgrade replaces contract bytecode and optional gas limit.
func (m *ContractManager) Upgrade(addr string, wasm []byte, gasLimit uint64) error {
	if len(wasm) == 0 {
		return ErrWASMRequired
	}
	c, ok := m.registry.Get(addr)
	if !ok {
		return fmt.Errorf("%w: %s", ErrContractNotFound, addr)
	}
	m.registry.mu.Lock()
	wasmCopy := make([]byte, len(wasm))
	copy(wasmCopy, wasm)
	c.WASM = wasmCopy
	if gasLimit > 0 {
		c.GasLimit = gasLimit
	}
	snapshot := cloneContract(c)
	m.registry.mu.Unlock()
	m.registry.persistContract(snapshot)
	m.registry.emit(ContractRegistryEvent{
		Type:     ContractRegistryEventUpgrade,
		Contract: snapshot,
		Caller:   snapshot.Owner,
		GasLimit: snapshot.GasLimit,
	})
	return nil
}

// Info returns contract metadata including owner and paused status.
func (m *ContractManager) Info(addr string) (*Contract, error) {
	c, ok := m.registry.Get(addr)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrContractNotFound, addr)
	}
	return cloneContract(c), nil
}
