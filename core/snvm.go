package core

import "fmt"

// SNVM represents the Synnergy Network Virtual Machine.
// Currently it operates as a very small stack machine used for executing
// arithmetic instructions contained within a Transaction's Program field.
// The structure is intentionally lightweight but can be extended with state or
// configuration as the VM evolves.
type SNVM struct{}

// NewSNVM creates a new virtual machine instance.
func NewSNVM() *SNVM { return &SNVM{} }

// Execute runs the opcodes associated with a transaction.  The function returns
// the top of the VM stack after execution or zero if the program produced no
// result.  Errors are returned for invalid programs such as stack underflows or
// division by zero.
func (vm *SNVM) Execute(tx *Transaction) (int64, error) {
	var stack []int64
	for _, ins := range tx.Program {
		switch ins.Op {
		case OpNoop:
			// do nothing
		case OpPush:
			stack = append(stack, ins.Value)
		case OpAdd, OpSub, OpMul, OpDiv:
			if len(stack) < 2 {
				return 0, fmt.Errorf("stack underflow")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var res int64
			switch ins.Op {
			case OpAdd:
				res = a + b
			case OpSub:
				res = a - b
			case OpMul:
				res = a * b
			case OpDiv:
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				res = a / b
			}
			stack = append(stack, res)
		default:
			return 0, fmt.Errorf("unknown opcode %d", ins.Op)
		}
	}

	if len(stack) == 0 {
		return 0, nil
	}
	return stack[len(stack)-1], nil
}
