package synnergy

import "testing"

func TestSNVMOpcodeLookup(t *testing.T) {
	code := SNVMOpcodeByName("gui_dex_screener_Liquidity")
	if code != 0x0004EE {
		t.Fatalf("unexpected opcode: %#x", code)
	}
}
