package core

// SNVM represents the Synnergy Network Virtual Machine.
type SNVM struct{}

// NewSNVM creates a new virtual machine instance.
func NewSNVM() *SNVM { return &SNVM{} }

// Execute runs the opcodes associated with a transaction.
func (vm *SNVM) Execute(tx *Transaction) error {
	// TODO: implement opcode execution
	return nil
}
