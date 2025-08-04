package core

// Basic opcode definitions used by the VM and gas table.
const (
	OpSub Opcode = iota
	OpMul
	OpDiv
	OpTransfer
)
