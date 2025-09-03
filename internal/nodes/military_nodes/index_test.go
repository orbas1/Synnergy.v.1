package militarynodes

import "testing"

func TestSimpleWarfareNodeLogistics(t *testing.T) {
	n := NewSimpleWarfareNode("m1")
	n.TrackLogistics("asset", "base", "ok")
	logs := n.Logistics()
	if len(logs) != 1 || logs[0].AssetID != "asset" {
		t.Fatalf("unexpected logs: %#v", logs)
	}
}
