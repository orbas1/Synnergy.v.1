package core

import "testing"

func TestSNVMArithmetic(t *testing.T) {
	vm := NewSNVM()
	tx := &Transaction{Program: []Instruction{
		{Op: OpPush, Value: 2},
		{Op: OpPush, Value: 3},
		{Op: OpAdd},
	}}
	res, err := vm.Execute(tx)
	if err != nil {
		t.Fatalf("execution error: %v", err)
	}
	if res != 5 {
		t.Fatalf("expected 5, got %d", res)
	}
}

func TestSNVMDivideByZero(t *testing.T) {
	vm := NewSNVM()
	tx := &Transaction{Program: []Instruction{
		{Op: OpPush, Value: 10},
		{Op: OpPush, Value: 0},
		{Op: OpDiv},
	}}
	if _, err := vm.Execute(tx); err == nil {
		t.Fatalf("expected division by zero error")
	}
}
