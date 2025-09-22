package tokens

import (
	"crypto/rand"
	"testing"
)

func TestSYN1100EncryptionWorkflow(t *testing.T) {
	token := NewSYN1100Token()
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("rand read: %v", err)
	}
	if err := token.SetEncryptionKey(key); err != nil {
		t.Fatalf("set key: %v", err)
	}

	if err := token.AddRecord(1, "alice", []byte("blood-type:O+")); err != nil {
		t.Fatalf("add record: %v", err)
	}

	if _, err := token.GetRecord(1, "bob"); err == nil {
		t.Fatal("expected access denial")
	}

	if err := token.GrantAccess(1, "bob"); err != nil {
		t.Fatalf("grant access: %v", err)
	}

	data, err := token.GetRecord(1, "bob")
	if err != nil {
		t.Fatalf("get record: %v", err)
	}
	if string(data) != "blood-type:O+" {
		t.Fatalf("unexpected data: %s", data)
	}

	if err := token.UpdateRecord(1, "bob", []byte("blood-type:A-")); err != nil {
		t.Fatalf("update record: %v", err)
	}

	plain, err := token.GetRecord(1, "alice")
	if err != nil {
		t.Fatalf("get updated: %v", err)
	}
	if string(plain) != "blood-type:A-" {
		t.Fatalf("unexpected updated data: %s", plain)
	}

	audit := token.AuditTrail(10)
	if len(audit) < 3 {
		t.Fatalf("expected audit entries, got %d", len(audit))
	}
}
