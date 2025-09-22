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
	latest, ok := n.Latest("subject")
	if !ok || latest.Latitude != 3.0 {
		t.Fatalf("unexpected latest record: %#v", latest)
	}
	summary, ok := n.Summary("subject")
	if !ok || summary.Count != 2 || summary.Bounds.MinLat != 1.0 || summary.Bounds.MaxLon != 4.0 {
		t.Fatalf("unexpected summary: %#v", summary)
	}
}

func TestGeospatialNodeValidationAndRetention(t *testing.T) {
	n := NewGeospatialNode(Address("g2"), WithGeospatialMaxHistory(2), WithGeospatialRetention(50*time.Millisecond))
	if err := n.Record("", 0, 0); err == nil {
		t.Fatalf("expected error for empty subject")
	}
	if err := n.Record("subject", 200, 0); err == nil {
		t.Fatalf("expected latitude validation error")
	}
	if err := n.Record("subject", 0, 200); err == nil {
		t.Fatalf("expected longitude validation error")
	}

	if err := n.Record("subject", 0, 0); err != nil {
		t.Fatalf("record: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	if err := n.Record("subject", 1, 1); err != nil {
		t.Fatalf("record: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	if err := n.Record("subject", 2, 2); err != nil {
		t.Fatalf("record: %v", err)
	}

	hist := n.History("subject")
	if len(hist) != 2 || hist[0].Latitude != 1 || hist[1].Latitude != 2 {
		t.Fatalf("unexpected retained history: %#v", hist)
	}
}

func TestGeospatialNodeHistoryWithin(t *testing.T) {
	n := NewGeospatialNode(Address("g3"))
	for i := 0; i < 5; i++ {
		if err := n.Record("subject", float64(i), float64(i)); err != nil {
			t.Fatalf("record: %v", err)
		}
		time.Sleep(time.Millisecond)
	}
	full := n.History("subject")
	if len(full) != 5 {
		t.Fatalf("expected 5 records got %d", len(full))
	}
	mid := full[1].Timestamp
	end := full[3].Timestamp
	records := n.HistoryWithin("subject", HistoryQuery{Since: mid, Until: end})
	if len(records) == 0 {
		t.Fatalf("expected filtered records")
	}
	for _, rec := range records {
		if rec.Timestamp.Before(mid) || rec.Timestamp.After(end) {
			t.Fatalf("record outside range: %#v", rec)
		}
	}
}
