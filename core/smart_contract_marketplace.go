package core

import (
	"context"
	"fmt"

	synn "synnergy"
	"synnergy/internal/telemetry"
)

// SmartContractMarketplace provides a minimal marketplace where users can deploy
// contracts and transfer ownership between parties. It wires the contract
// registry and manager together so higher level interfaces (CLI/GUI) only need a
// single component. The marketplace is concurrency-safe through the underlying
// registry mutexes and relies on gas pricing to prevent abuse.
type SmartContractMarketplace struct {
	registry *ContractRegistry
	manager  *ContractManager
}

// NewSmartContractMarketplace initialises a marketplace backed by the provided
// virtual machine. The VM is started if necessary so contracts can be executed
// immediately after deployment.
func NewSmartContractMarketplace(vm VirtualMachine) *SmartContractMarketplace {
	if !vm.Status() {
		_ = vm.Start()
	}
	reg := NewContractRegistry(vm)
	return &SmartContractMarketplace{registry: reg, manager: NewContractManager(reg)}
}

// DeployContract compiles and stores the given WASM contract. A manifest and gas
// limit may be supplied; the caller must also provide an owner address. The
// deployment fails if the supplied gas is below the registered cost or if the
// contract already exists.
func (m *SmartContractMarketplace) DeployContract(ctx context.Context, wasm []byte, manifest string, gasLimit uint64, owner string) (string, error) {
	ctx, span := telemetry.Tracer().Start(ctx, "SmartContractMarketplace.DeployContract")
	defer span.End()

	required := synn.GasCost("DeploySmartContract")
	if gasLimit < required {
		return "", fmt.Errorf("%w: need %d", ErrInsufficientGas, required)
	}
	return m.registry.Deploy(wasm, manifest, gasLimit, owner)
}

// TradeContract transfers ownership of an existing contract. It returns an
// error if the contract cannot be found or the underlying manager rejects the
// transfer.
func (m *SmartContractMarketplace) TradeContract(ctx context.Context, addr, newOwner string) error {
	ctx, span := telemetry.Tracer().Start(ctx, "SmartContractMarketplace.TradeContract")
	defer span.End()
	return m.manager.Transfer(ctx, addr, newOwner)
}
