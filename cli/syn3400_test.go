package cli

import (
        "strings"
        "testing"

        "synnergy/internal/tokens"
)

func TestSyn3400Lifecycle(t *testing.T) {
        forexRegistry = tokens.NewForexRegistry()
        out, err := execCommand("syn3400", "register", "--base", "USD", "--quote", "EUR", "--rate", "1.2")
        if err != nil {
                t.Fatalf("register: %v", err)
        }
        id := strings.TrimSpace(out)
        if id == "" {
                t.Fatalf("expected id, got %q", out)
        }
        if _, err := execCommand("syn3400", "update", id, "1.3"); err != nil {
                t.Fatalf("update: %v", err)
        }
        out, err = execCommand("syn3400", "get", id)
        if err != nil {
                t.Fatalf("get: %v", err)
        }
        if !strings.Contains(out, "Rate:1.300000") {
                t.Fatalf("unexpected output: %s", out)
        }
}

func TestSyn3400InvalidRegister(t *testing.T) {
        forexRegistry = tokens.NewForexRegistry()
        if _, err := execCommand("syn3400", "register", "--base", "", "--quote", "EUR", "--rate", "1.2"); err == nil {
                t.Fatal("expected error for missing base")
        }
        if _, err := execCommand("syn3400", "register", "--base", "USD", "--quote", "EUR", "--rate", "-1"); err == nil {
                t.Fatal("expected error for negative rate")
        }
}
