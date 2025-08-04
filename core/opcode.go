package core

// Opcode defines the set of operations understood by the SNVM.
//
// The opcodes below provide a very small but functional instruction set for
// demonstration and testing of the virtual machine.  They intentionally mirror
// common arithmetic instructions one might find in a typical stack based VM.
// Additional instructions can be added over time as the SNVM grows.
type Opcode byte

const (
	OpNoop Opcode = iota // no operation

	// Stack / arithmetic operations
	OpPush // push literal onto the stack
	OpAdd  // pop two values, push a + b
	OpSub  // pop two values, push a - b
	OpMul  // pop two values, push a * b
	OpDiv  // pop two values, push a / b

	// Application level operation used elsewhere in the codebase.  It is
	// retained to maintain compatibility with existing modules such as the
	// gas accounting logic.
	OpTransfer
)
