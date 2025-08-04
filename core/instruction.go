package core

// Instruction represents a single operation and optional value for the
// Synnergy Virtual Machine.  Only OpPush utilises the Value field.
// Other opcodes operate purely on the VM stack.
type Instruction struct {
	Op    Opcode
	Value int64
}
