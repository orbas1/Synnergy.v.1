package core

// Basic opcodes used by the VM and tests. In the full system these values are
// generated, but for the purposes of the core library they are defined
// explicitly.
const (
	OpNoop Opcode = iota + 1
	OpPush
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpTransfer
)
