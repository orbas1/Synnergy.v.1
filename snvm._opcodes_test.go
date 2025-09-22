package synnergy

import "testing"

func TestSNVMOpcodeLookup(t *testing.T) {
	t.Parallel()
	if code := SNVMOpcodeByName("gui_dex_screener_liquidity"); code != 0x0004EE {
		t.Fatalf("unexpected opcode: %#x", code)
	}
	if op, ok := SNVMOpcodeByCode(0x0004EF); !ok || op.Name != "gui_security_operations_RaiseAlert" {
		t.Fatalf("unexpected metadata: %#v %v", op, ok)
	}
	if _, ok := SNVMOpcodeByCode(0xFFFFFF); ok {
		t.Fatalf("expected missing opcode")
	}
}

func TestSNVMOpcodeCatalogueIsolation(t *testing.T) {
	t.Parallel()
	cat := SNVMOpcodeCatalogue()
	if len(cat) == 0 {
		t.Fatal("catalogue should not be empty")
	}
	cat[0].Name = "tampered"
	other := SNVMOpcodeCatalogue()
	if other[0].Name == "tampered" {
		t.Fatal("catalogue must be returned as a copy")
	}
}
