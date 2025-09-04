package synnergy

import "testing"

func TestSNVMOpcodeLookup(t *testing.T) {
	if code := SNVMOpcodeByName("gui_dex_screener_Liquidity"); code != 0x0004EE {
		t.Fatalf("unexpected opcode: %#x", code)
	}
	if code := SNVMOpcodeByName("gui_security_operations_RaiseAlert"); code != 0x0004EF {
		t.Fatalf("unexpected opcode: %#x", code)
	}
	if code := SNVMOpcodeByName("gui_smart_contract_marketplace_ListContract"); code != 0x0004F0 {
		t.Fatalf("unexpected opcode: %#x", code)
	}
}
