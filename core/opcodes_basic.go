package core

// Basic opcode definitions used by the VM and gas table.
const (
	OpNoop Opcode = iota
	OpPush
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpTransfer
)
