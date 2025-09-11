package core

// Basic arithmetic opcodes for the simple SNVM stack machine.
const (
	OpNoop Opcode = iota
	OpPush
	OpAdd
	OpSub
	OpMul
	OpDiv
	// OpMod returns the remainder of the division of the top two stack values.
	OpMod
)
