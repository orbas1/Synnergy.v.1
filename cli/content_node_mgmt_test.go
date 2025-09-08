package cli

import (
	"strings"
	"testing"
)

func TestContentNodeCLI(t *testing.T) {
	if _, err := execCommand("content_node", "register", "id1", "name1"); err != nil {
		t.Fatalf("register: %v", err)
	}
	out, err := execCommand("content_node", "--json", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "id1") {
		t.Fatalf("expected id1 in list, got %s", out)
	}
	if _, err := execCommand("content_node", "unregister", "id1"); err != nil {
		t.Fatalf("unregister: %v", err)
	}
}
