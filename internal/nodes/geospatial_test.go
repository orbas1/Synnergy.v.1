package nodes

import (
	"testing"
	"time"
)

func TestGeospatialNodeRecordHistory(t *testing.T) {
	n := NewGeospatialNode(Address("g1"))
	if err := n.Record("subject", 1.0, 2.0); err != nil {
		t.Fatalf("record: %v", err)
	}
	time.Sleep(time.Millisecond)
	if err := n.Record("subject", 3.0, 4.0); err != nil {
		t.Fatalf("record: %v", err)
	}
	hist := n.History("subject")
	if len(hist) != 2 {
		t.Fatalf("expected 2 records got %d", len(hist))
	}
	if hist[0].Latitude != 1.0 || hist[1].Latitude != 3.0 {
		t.Fatalf("unexpected history ordering: %#v", hist)
	}
}
