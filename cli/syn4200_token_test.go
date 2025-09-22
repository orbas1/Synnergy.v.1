package cli

import "testing"

func TestSyn4200TokenLifecycle(t *testing.T) {
	if _, err := execCommand("syn4200_token", "donate", "AAA", "--from", "alice", "--amt", "50", "--purpose", "school"); err != nil {
		t.Fatalf("donate: %v", err)
	}
	out, err := execCommand("syn4200_token", "progress", "AAA")
	if err != nil {
		t.Fatalf("progress: %v", err)
	}
	if out != "50" {
		t.Fatalf("expected 50, got %s", out)
	}
}

func TestSyn4200TokenValidation(t *testing.T) {
	if _, err := execCommand("syn4200_token", "donate", "AAA", "--from", "", "--amt", "10"); err == nil {
		t.Fatal("expected error for missing from")
	}
	if _, err := execCommand("syn4200_token", "progress", "BBB"); err == nil {
		t.Fatal("expected error for unknown campaign")
	}
}
