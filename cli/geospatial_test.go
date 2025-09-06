package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestGeospatialRecordHistory verifies geospatial records can be stored and retrieved.
func TestGeospatialRecordHistory(t *testing.T) {
	if _, err := execCommand("geospatial", "record", "subj", "1.23", "2.34", "--json"); err != nil {
		t.Fatalf("record: %v", err)
	}
	out, err := execCommand("geospatial", "history", "subj", "--json")
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var recs []struct {
		Subject string  `json:"subject"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	}
	if err := json.Unmarshal([]byte(out), &recs); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(recs) != 1 || recs[0].Subject != "subj" || recs[0].Lat != 1.23 || recs[0].Lon != 2.34 {
		t.Fatalf("unexpected record: %v", recs)
	}
}
