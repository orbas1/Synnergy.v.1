package core

import "errors"

// Instruction represents a single operation and optional value for the
// Synnergy Virtual Machine. Only OpPush utilises the Value field. Other
// opcodes operate purely on the VM stack.
type Instruction struct {
	Op    Opcode
	Value int64
}

// Validate ensures the instruction conforms to VM rules.
// For non-push opcodes the Value must be zero.
func (i Instruction) Validate() error {
	if i.Op != OpPush && i.Value != 0 {
		return errors.New("value only permitted with OpPush")
	}
	return nil
}
