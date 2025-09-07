package cli

import "testing"

func TestSyn500Lifecycle(t *testing.T) {
	syn500Token = nil
	out, err := execCommand("syn500", "create", "--name", "Loyalty", "--symbol", "LOY", "--owner", "alice", "--dec", "2", "--supply", "10")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if out != "token created" {
		t.Fatalf("unexpected output: %s", out)
	}
	if _, err := execCommand("syn500", "grant", "bob", "--tier", "1", "--max", "2"); err != nil {
		t.Fatalf("grant: %v", err)
	}
	if _, err := execCommand("syn500", "use", "bob"); err != nil {
		t.Fatalf("use1: %v", err)
	}
	if _, err := execCommand("syn500", "use", "bob"); err != nil {
		t.Fatalf("use2: %v", err)
	}
	if _, err := execCommand("syn500", "use", "bob"); err == nil {
		t.Fatal("expected usage limit error")
	}
}

func TestSyn500Validation(t *testing.T) {
	syn500Token = nil
	if _, err := execCommand("syn500", "create", "--name", "", "--symbol", "LOY", "--owner", "alice", "--dec", "1", "--supply", "10"); err == nil {
		t.Fatal("expected error for name")
	}
	if _, err := execCommand("syn500", "create", "--name", "Loy", "--symbol", "LOY", "--owner", "", "--dec", "1", "--supply", "10"); err == nil {
		t.Fatal("expected error for owner")
	}
	if _, err := execCommand("syn500", "create", "--name", "Loy", "--symbol", "LOY", "--owner", "alice", "--dec", "0", "--supply", "10"); err == nil {
		t.Fatal("expected error for decimals")
	}
	if _, err := execCommand("syn500", "create", "--name", "Loy", "--symbol", "LOY", "--owner", "alice", "--dec", "1", "--supply", "0"); err == nil {
		t.Fatal("expected error for supply")
	}
	if _, err := execCommand("syn500", "create", "--name", "Loy", "--symbol", "LOY", "--owner", "alice", "--dec", "1", "--supply", "5"); err != nil {
		t.Fatalf("create valid: %v", err)
	}
	if _, err := execCommand("syn500", "grant", "bob", "--tier", "0", "--max", "1"); err == nil {
		t.Fatal("expected error for tier")
	}
	if _, err := execCommand("syn500", "grant", "bob", "--tier", "1", "--max", "0"); err == nil {
		t.Fatal("expected error for max")
	}
}
