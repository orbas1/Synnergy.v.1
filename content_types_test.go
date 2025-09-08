package synnergy

import (
	"strings"
	"testing"
	"time"
)

func TestNewContentMetaFromData(t *testing.T) {
	data := []byte("hello world")
	cm := NewContentMetaFromData("id1", "greeting", data)

	if cm.ID != "id1" || cm.Name != "greeting" {
		t.Fatalf("unexpected metadata fields: %+v", cm)
	}
	if cm.Size != int64(len(data)) {
		t.Fatalf("size mismatch: got %d want %d", cm.Size, len(data))
	}
	if len(cm.Hash) != 64 {
		t.Fatalf("hash length = %d, want 64", len(cm.Hash))
	}
	if time.Since(cm.Created) > time.Second {
		t.Fatalf("creation time too old: %v", cm.Created)
	}
	if err := cm.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestContentMetaValidateErrors(t *testing.T) {
	cases := []ContentMeta{
		{Name: "name", Size: 1, Hash: strings.Repeat("0", 64)},            // missing ID
		{ID: "id", Size: 1, Hash: strings.Repeat("0", 64)},                // missing name
		{ID: "id", Name: "name", Size: -1, Hash: strings.Repeat("0", 64)}, // negative size
		{ID: "id", Name: "name", Size: 1, Hash: "deadbeef"},               // short hash
	}

	for i, cm := range cases {
		if err := cm.Validate(); err == nil {
			t.Fatalf("case %d expected error", i)
		}
	}
}
