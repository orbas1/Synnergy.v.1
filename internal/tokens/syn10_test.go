package tokens

import "testing"

func TestSYN10ComplianceFlow(t *testing.T) {
	tok := NewSYN10Token(1, "CBDC", "CBDC", "Synnergy Reserve", 1.0, 2)
	tok.Mint("treasury", 1_000)

	tok.RegisterParticipant("treasury", TierInstitution, 500, map[string]string{"role": "issuer"})
	tok.RegisterParticipant("alice", TierRetail, 500, map[string]string{"kyc": "level2"})
	tok.RegisterParticipant("bob", TierRetail, 250, map[string]string{"kyc": "level1"})

	if err := tok.TransferCBDC("treasury", "alice", 400, map[string]string{"reason": "payout"}); err != nil {
		t.Fatalf("transfer: %v", err)
	}

	if err := tok.TransferCBDC("treasury", "alice", 200, nil); err != ErrCBDCDailyLimitExceeded {
		t.Fatalf("expected limit exceeded, got %v", err)
	}

	if err := tok.FlagParticipant("alice", true, map[string]string{"case": "aml-review"}); err != nil {
		t.Fatalf("flag: %v", err)
	}
	if err := tok.TransferCBDC("treasury", "alice", 1, nil); err != ErrCBDCParticipantFlagged {
		t.Fatalf("expected flagged error, got %v", err)
	}

	if err := tok.FlagParticipant("alice", false, nil); err != nil {
		t.Fatalf("unflag: %v", err)
	}

	tok.SetCircuitBreaker(true)
	if err := tok.TransferCBDC("treasury", "bob", 10, nil); err != ErrCBDCCircuitBreakerActive {
		t.Fatalf("expected circuit breaker error, got %v", err)
	}
	tok.SetCircuitBreaker(false)

	if err := tok.TransferCBDC("treasury", "bob", 10, map[string]string{"reason": "stipend"}); err != nil {
		t.Fatalf("transfer after reset: %v", err)
	}

	audits := tok.AuditTrail(10)
	if len(audits) != 2 {
		t.Fatalf("expected 2 audit entries, got %d", len(audits))
	}
	if audits[len(audits)-1].To != "bob" {
		t.Fatalf("unexpected audit recipient: %+v", audits[len(audits)-1])
	}
}
