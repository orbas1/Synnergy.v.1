package tokens

import "testing"

func TestItemRegistryLifecycle(t *testing.T) {
	reg := NewItemRegistry()
	item := reg.CreateItem("alice", "Rare Sword", "Legendary weapon", map[string]string{"rarity": "legendary"})
	if item.Owner != "alice" {
		t.Fatalf("unexpected owner: %s", item.Owner)
	}

	if err := reg.UpdateAttributes(item.ItemID, map[string]string{"power": "99"}); err != nil {
		t.Fatalf("update attributes: %v", err)
	}
	if err := reg.TransferItem(item.ItemID, "bob"); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if err := reg.ArchiveItem(item.ItemID, "bob"); err != nil {
		t.Fatalf("archive: %v", err)
	}

	snapshot, ok := reg.GetItem(item.ItemID)
	if !ok || !snapshot.Archived || snapshot.Owner != "bob" {
		t.Fatalf("unexpected snapshot: %+v", snapshot)
	}

	events, ok := reg.Events(item.ItemID, 10)
	if !ok || len(events) != 4 {
		t.Fatalf("unexpected events: %+v", events)
	}
	if events[len(events)-1].Type != "archive" {
		t.Fatalf("expected archive event, got %+v", events[len(events)-1])
	}
}
