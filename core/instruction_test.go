package core

import "testing"

func TestInstructionValidate(t *testing.T) {
	ins := Instruction{Op: OpAdd, Value: 1}
	if err := ins.Validate(); err == nil {
		t.Fatalf("expected error for value with non-push op")
	}
	ins2 := Instruction{Op: OpPush, Value: 5}
	if err := ins2.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
