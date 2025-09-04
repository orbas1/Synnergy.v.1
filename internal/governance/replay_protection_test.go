package governance

import "testing"

func TestReplayProtectorSeen(t *testing.T) {
	r := NewReplayProtector()
	if r.Seen("abc") {
		t.Fatal("ID should not be seen first time")
	}
	if !r.Seen("abc") {
		t.Fatal("ID should be seen second time")
	}
}
