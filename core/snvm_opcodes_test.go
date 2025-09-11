package core

import "testing"

// TestOpcodeValues ensures the SNVM opcode constants remain stable.
func TestOpcodeValues(t *testing.T) {
	cases := []struct {
		op   Opcode
		want Opcode
	}{
		{OpNoop, 0},
		{OpPush, 1},
		{OpAdd, 2},
		{OpSub, 3},
		{OpMul, 4},
		{OpDiv, 5},
		{OpMod, 6},
	}
	for _, tc := range cases {
		if tc.op != tc.want {
			t.Fatalf("opcode %v expected %d", tc.op, tc.want)
		}
	}
}
