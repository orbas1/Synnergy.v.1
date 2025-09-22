package core

import (
	"context"

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

	c, ok := m.registry.Get(addr)
	if !ok {
		return ierr.New(ierr.NotFound, "contract not found")
	}
	m.registry.mu.Lock()
	c.Owner = newOwner
	m.registry.mu.Unlock()
	return nil
}

// Pause disables contract execution.
func (m *ContractManager) Pause(ctx context.Context, addr string) error {
	ctx, span := telemetry.Tracer("core.contracts").Start(ctx, "ContractManager.Pause")
	defer span.End()

	c, ok := m.registry.Get(addr)
	if !ok {
		return ierr.New(ierr.NotFound, "contract not found")
	}
	m.registry.mu.Lock()
	c.Paused = true
	m.registry.mu.Unlock()
	return nil
}

// Resume enables execution for a paused contract.
func (m *ContractManager) Resume(ctx context.Context, addr string) error {
	ctx, span := telemetry.Tracer("core.contracts").Start(ctx, "ContractManager.Resume")
	defer span.End()

	c, ok := m.registry.Get(addr)
	if !ok {
		return ierr.New(ierr.NotFound, "contract not found")
	}
	m.registry.mu.Lock()
	c.Paused = false
	m.registry.mu.Unlock()
	return nil
}

// Upgrade replaces contract bytecode and optional gas limit.
func (m *ContractManager) Upgrade(ctx context.Context, addr string, wasm []byte, gasLimit uint64) error {
	ctx, span := telemetry.Tracer("core.contracts").Start(ctx, "ContractManager.Upgrade")
	defer span.End()

	if len(wasm) == 0 {
		return ierr.New(ierr.Invalid, "wasm bytecode required")
	}
	c, ok := m.registry.Get(addr)
	if !ok {
		return ierr.New(ierr.NotFound, "contract not found")
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
func (m *ContractManager) Info(ctx context.Context, addr string) (*Contract, error) {
	ctx, span := telemetry.Tracer("core.contracts").Start(ctx, "ContractManager.Info")
	defer span.End()

	c, ok := m.registry.Get(addr)
	if !ok {
		return nil, ierr.New(ierr.NotFound, "contract not found")
	}
	return c, nil
}
