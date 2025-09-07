package cli

import "testing"

func TestSyn3700Lifecycle(t *testing.T) {
	syn3700 = nil
	if _, err := execCommand("syn3700", "init", "--name", "Index", "--symbol", "IDX"); err != nil {
		t.Fatalf("init: %v", err)
	}
	if _, err := execCommand("syn3700", "add", "AAA", "0.5"); err != nil {
		t.Fatalf("add1: %v", err)
	}
	if _, err := execCommand("syn3700", "add", "BBB", "1.5"); err != nil {
		t.Fatalf("add2: %v", err)
	}
	out, err := execCommand("syn3700", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	expected := "AAA 0.50\nBBB 1.50"
	if out != expected {
		t.Fatalf("unexpected list: %q", out)
	}
	out, err = execCommand("syn3700", "value", "AAA:2", "BBB:3")
	if err != nil {
		t.Fatalf("value: %v", err)
	}
	if out != "5.50" {
		t.Fatalf("expected 5.50, got %s", out)
	}
	if _, err := execCommand("syn3700", "remove", "AAA"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	out, err = execCommand("syn3700", "list")
	if err != nil {
		t.Fatalf("list2: %v", err)
	}
	if out != "BBB 1.50" {
		t.Fatalf("unexpected list2: %q", out)
	}
}

func TestSyn3700Validation(t *testing.T) {
	syn3700 = nil
	if _, err := execCommand("syn3700", "init", "--name", "", "--symbol", "IDX"); err == nil {
		t.Fatal("expected error for missing name")
	}
	if _, err := execCommand("syn3700", "init", "--name", "Index", "--symbol", "IDX"); err != nil {
		t.Fatalf("init: %v", err)
	}
	if _, err := execCommand("syn3700", "add", "", "1"); err == nil {
		t.Fatal("expected error for empty token symbol")
	}
	if _, err := execCommand("syn3700", "add", "AAA", "-1"); err == nil {
		t.Fatal("expected error for negative weight")
	}
	if _, err := execCommand("syn3700", "value", "badpair"); err == nil {
		t.Fatal("expected error for malformed price pair")
	}
}
