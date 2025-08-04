package core

import "errors"

// ContractManager provides administrative operations over deployed contracts.
type ContractManager struct {
	registry *ContractRegistry
}

// NewContractManager wires a manager to an existing registry.
func NewContractManager(reg *ContractRegistry) *ContractManager {
	return &ContractManager{registry: reg}
}

// Transfer changes the owner of a contract.
func (m *ContractManager) Transfer(addr, newOwner string) error {
	c, ok := m.registry.Get(addr)
	if !ok {
		return errors.New("contract not found")
	}
	m.registry.mu.Lock()
	c.Owner = newOwner
	m.registry.mu.Unlock()
	return nil
}

// Pause disables contract execution.
func (m *ContractManager) Pause(addr string) error {
	c, ok := m.registry.Get(addr)
	if !ok {
		return errors.New("contract not found")
	}
	m.registry.mu.Lock()
	c.Paused = true
	m.registry.mu.Unlock()
	return nil
}

// Resume enables execution for a paused contract.
func (m *ContractManager) Resume(addr string) error {
	c, ok := m.registry.Get(addr)
	if !ok {
		return errors.New("contract not found")
	}
	m.registry.mu.Lock()
	c.Paused = false
	m.registry.mu.Unlock()
	return nil
}

// Upgrade replaces contract bytecode and optional gas limit.
func (m *ContractManager) Upgrade(addr string, wasm []byte, gasLimit uint64) error {
	if len(wasm) == 0 {
		return errors.New("wasm bytecode required")
	}
	c, ok := m.registry.Get(addr)
	if !ok {
		return errors.New("contract not found")
	}
	m.registry.mu.Lock()
	c.WASM = wasm
	if gasLimit > 0 {
		c.GasLimit = gasLimit
	}
	m.registry.mu.Unlock()
	return nil
}

// Info returns contract metadata including owner and paused status.
func (m *ContractManager) Info(addr string) (*Contract, error) {
	c, ok := m.registry.Get(addr)
	if !ok {
		return nil, errors.New("contract not found")
	}
	return c, nil
}
