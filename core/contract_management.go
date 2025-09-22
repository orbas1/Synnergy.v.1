package core

import (
	"context"
	"fmt"
	"strings"

	ierr "synnergy/internal/errors"
	"synnergy/internal/telemetry"
)

// ContractManager provides administrative operations over deployed contracts.
type ContractManager struct {
	registry *ContractRegistry
}

// NewContractManager wires a manager to an existing registry.
func NewContractManager(reg *ContractRegistry) *ContractManager {
	return &ContractManager{registry: reg}
}

// Transfer changes the owner of a contract.
func (m *ContractManager) Transfer(ctx context.Context, addr, newOwner string) error {
	ctx, span := telemetry.Tracer("core.contracts").Start(ctx, "ContractManager.Transfer")
	defer span.End()

	newOwner = strings.TrimSpace(newOwner)
	if newOwner == "" {
		return ierr.New(ierr.Invalid, "new contract owner cannot be empty")
	}
	c, ok := m.registry.Get(addr)
	if !ok {
		return ierr.New(ierr.NotFound, fmt.Sprintf("contract not found: %s", addr))
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
func (m *ContractManager) Pause(ctx context.Context, addr string) error {
	ctx, span := telemetry.Tracer("core.contracts").Start(ctx, "ContractManager.Pause")
	defer span.End()

	c, ok := m.registry.Get(addr)
	if !ok {
		return ierr.New(ierr.NotFound, fmt.Sprintf("contract not found: %s", addr))
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
func (m *ContractManager) Resume(ctx context.Context, addr string) error {
	ctx, span := telemetry.Tracer("core.contracts").Start(ctx, "ContractManager.Resume")
	defer span.End()

	c, ok := m.registry.Get(addr)
	if !ok {
		return ierr.New(ierr.NotFound, fmt.Sprintf("contract not found: %s", addr))
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
func (m *ContractManager) Upgrade(ctx context.Context, addr string, wasm []byte, gasLimit uint64) error {
	ctx, span := telemetry.Tracer("core.contracts").Start(ctx, "ContractManager.Upgrade")
	defer span.End()

	if len(wasm) == 0 {
		return ierr.New(ierr.Invalid, ErrWASMRequired.Error())
	}
	c, ok := m.registry.Get(addr)
	if !ok {
		return ierr.New(ierr.NotFound, fmt.Sprintf("contract not found: %s", addr))
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
func (m *ContractManager) Info(ctx context.Context, addr string) (*Contract, error) {
	ctx, span := telemetry.Tracer("core.contracts").Start(ctx, "ContractManager.Info")
	defer span.End()

	c, ok := m.registry.Get(addr)
	if !ok {
		return nil, ierr.New(ierr.NotFound, fmt.Sprintf("contract not found: %s", addr))
	}
	return cloneContract(c), nil
}
