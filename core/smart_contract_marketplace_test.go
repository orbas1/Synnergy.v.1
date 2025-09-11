package core

import (
	"context"
	"testing"

	synn "synnergy"
)

// TestSmartContractMarketplaceDeployAndTrade ensures contracts can be deployed
// and ownership transferred via the marketplace with proper gas accounting.
func TestSmartContractMarketplaceDeployAndTrade(t *testing.T) {
	vm := NewSimpleVM()
	if err := vm.Start(); err != nil {
		t.Fatalf("vm start: %v", err)
	}
	m := NewSmartContractMarketplace(vm)

	gas := synn.GasCost("DeploySmartContract")
	addr, err := m.DeployContract(context.Background(), []byte{0x00}, "", gas, "alice")
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	if addr == "" {
		t.Fatal("expected contract address")
	}

	tradeGas := synn.GasCost("TradeContract")
	if err := m.TradeContract(context.Background(), addr, "bob", tradeGas); err != nil {
		t.Fatalf("trade: %v", err)
	}
	c, ok := m.registry.Get(addr)
	if !ok || c.Owner != "bob" {
		t.Fatalf("expected owner bob, got %+v", c)
	}
	// insufficient gas
	if err := m.TradeContract(context.Background(), addr, "carol", tradeGas-1); err == nil {
		t.Fatalf("expected insufficient gas error")
	}
}
