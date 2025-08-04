package core

// Opcode defines the set of operations understood by the SNVM.
type Opcode byte

const (
	OpNoop Opcode = iota
	OpTransfer
)
