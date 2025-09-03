package smartcontracts

import (
	"os"
	"testing"
)

// TestTemplatesExist ensures that the predefined WASM templates are present.
func TestTemplatesExist(t *testing.T) {
        templates := []string{
                "token_faucet.wasm",
                "storage_market.wasm",
                "dao_governance.wasm",
                "nft_minting.wasm",
                "ai_model_market.wasm",
                "escrow_payment.wasm",
                "cross_chain_bridge.wasm",
                "multisig_wallet.wasm",
                "regulatory_compliance.wasm",
        }
        for _, tpl := range templates {
                if _, err := os.Stat(tpl); err != nil {
                        t.Fatalf("template %s not found: %v", tpl, err)
                }
        }
}
