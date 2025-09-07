package cli

import (
	"encoding/json"
	"testing"

	"synnergy/core"
)

func TestSyn3900Lifecycle(t *testing.T) {
	benefitRegistry = core.NewBenefitRegistry()
	out, err := execCommand("syn3900", "register", "bob", "housing", "200")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if out != "1" {
		t.Fatalf("expected id 1, got %s", out)
	}
	if _, err := execCommand("syn3900", "claim", "1"); err != nil {
		t.Fatalf("claim: %v", err)
	}
	out, err = execCommand("syn3900", "get", "1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	var b struct{ Claimed bool }
	if err := json.Unmarshal([]byte(out), &b); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !b.Claimed {
		t.Fatalf("expected claimed=true")
	}
}

func TestSyn3900Validation(t *testing.T) {
	benefitRegistry = core.NewBenefitRegistry()
	if _, err := execCommand("syn3900", "register", "", "prog", "10"); err == nil {
		t.Fatal("expected error for missing recipient")
	}
	if _, err := execCommand("syn3900", "register", "alice", "prog", "0"); err == nil {
		t.Fatal("expected error for amount")
	}
	if _, err := execCommand("syn3900", "claim", "1"); err == nil {
		t.Fatal("expected error for unknown id")
	}
}
