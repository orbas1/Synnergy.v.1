package tokens

import "testing"

func TestSYN20Controls(t *testing.T) {
	tok := NewSYN20Token(20, "Permissioned", "S20", 2)
	tok.Mint("issuer", 1000)

	tok.AddToWhitelist("alice")
	tok.SetWhitelistRequirement(true)
	tok.SetDailyLimit("issuer", 500)

	if err := tok.Transfer("issuer", "alice", 400); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	tok.RecordAudit("issuer", "alice", 400, "initial grant")

	if err := tok.Transfer("issuer", "alice", 200); err != ErrSYN20LimitExceeded {
		t.Fatalf("expected limit exceeded, got %v", err)
	}

	tok.Freeze("alice")
	if err := tok.Transfer("issuer", "alice", 10); err == nil {
		t.Fatal("expected frozen address error")
	}
	tok.Unfreeze("alice")

	tok.SetCircuitBreaker(true)
	if err := tok.Transfer("issuer", "alice", 10); err != ErrSYN20CircuitBreakerActive {
		t.Fatalf("expected circuit breaker error, got %v", err)
	}
	tok.SetCircuitBreaker(false)

	if err := tok.Transfer("issuer", "alice", 10); err != nil {
		t.Fatalf("unexpected transfer error: %v", err)
	}

	audits := tok.AuditTrail(5)
	if len(audits) == 0 {
		t.Fatal("expected audit entry")
	}

	tok.RemoveFromWhitelist("alice")
	if err := tok.Transfer("issuer", "alice", 1); err != ErrSYN20WhitelistViolation {
		t.Fatalf("expected whitelist violation, got %v", err)
	}
}
