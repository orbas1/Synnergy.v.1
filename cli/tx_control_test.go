package cli

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestTxControlSchedule(t *testing.T) {
	exec := strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10)
	out, err := execCommand("tx", "control", "schedule", "a", "b", "1", "1", "0", exec)
	if err != nil {
		t.Fatalf("schedule failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "txID") {
		t.Fatalf("unexpected output: %s", out)
	}
}
