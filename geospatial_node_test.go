package synnergy

import "testing"

func TestGeospatialNodeRecordAndHistory(t *testing.T) {
	n := NewGeospatialNode()
	n.Record("asset1", 1.23, 4.56)
	hist := n.History("asset1")
	if len(hist) != 1 {
		t.Fatalf("expected 1 record, got %d", len(hist))
	}
	if hist[0].Latitude != 1.23 || hist[0].Longitude != 4.56 {
		t.Fatalf("unexpected record %#v", hist[0])
	}
	// ensure returned slice is a copy
	hist[0].Latitude = 9.0
	if n.History("asset1")[0].Latitude == 9.0 {
		t.Fatalf("modifying history should not affect stored data")
	}
}
