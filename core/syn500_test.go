package core

import (
	"testing"
	"time"
)

func TestSyn500Usage(t *testing.T) {
	tok := NewSYN500Token("Loyalty", "LOY", "alice", 2, 10)
	tok.Grant("bob", 1, 2, time.Millisecond)
	if err := tok.Use("bob"); err != nil {
		t.Fatalf("use1: %v", err)
	}
	if err := tok.Use("bob"); err != nil {
		t.Fatalf("use2: %v", err)
	}
	if err := tok.Use("bob"); err == nil {
		t.Fatalf("expected usage limit error")
	}
	time.Sleep(2 * time.Millisecond)
	if err := tok.Use("bob"); err != nil {
		t.Fatalf("expected reset usage: %v", err)
	}
	st, ok := tok.Status("bob")
	if !ok || st.Used == 0 {
		t.Fatalf("expected usage state: %+v", st)
	}
	tele := tok.Telemetry()
	if tele.Grants != 1 {
		t.Fatalf("unexpected telemetry: %+v", tele)
	}
}
