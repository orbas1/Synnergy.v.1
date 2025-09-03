package errors

import (
	"fmt"
	"testing"
)

func TestErrorWrapping(t *testing.T) {
	base := fmt.Errorf("boom")
	err := Wrap(Invalid, "failed", base)
	if err.Error() == "" {
		t.Fatalf("empty error")
	}
	if !IsCode(err, Invalid) {
		t.Fatalf("code mismatch")
	}
	if err.Unwrap() != base {
		t.Fatalf("unwrap mismatch")
	}
}
