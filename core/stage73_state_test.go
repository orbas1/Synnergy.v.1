package core

import "testing"

func TestStage73StateRoundTrip(t *testing.T) {
	tok := NewSYN3700Token("Index", "IDX")
	_ = tok.AddController("ctrl", "ctrl")
	_ = tok.AddComponent("AAA", 0.5, 0.1, "ctrl")
	grants := NewGrantRegistry()
	grantID, _ := grants.CreateGrant("alice", "grant", 10, "ctrl")
	benefits := NewBenefitRegistry()
	benefitID, _ := benefits.RegisterBenefit("alice", "housing", 5, "ctrl")
	state := CaptureStage73State(tok, grants, benefits)

	tok2 := NewSYN3700Token("", "")
	grants2 := NewGrantRegistry()
	benefits2 := NewBenefitRegistry()
	ApplyStage73State(state, tok2, grants2, benefits2)

	if tele := tok2.Telemetry(); tele.ComponentCount != 1 {
		t.Fatalf("telemetry mismatch: %+v", tele)
	}
	if _, ok := grants2.GetGrant(grantID); !ok {
		t.Fatalf("grant not restored")
	}
	if _, ok := benefits2.GetBenefit(benefitID); !ok {
		t.Fatalf("benefit not restored")
	}
}
