package core

import "testing"

// TestSYN4200Token verifies donation handling, campaign tracking, and data isolation.
func TestSYN4200Token(t *testing.T) {
	token := NewSYN4200Token()

	// Initial donation should create the campaign.
	token.Donate("CHAR", "alice", 100, "school")
	// Second donation to same campaign from another donor.
	token.Donate("CHAR", "bob", 50, "ignored purpose")
	// Donation to a separate campaign.
	token.Donate("WATER", "carol", 200, "wells")

	// Verify progress for existing campaigns.
	if amt, ok := token.CampaignProgress("CHAR"); !ok || amt != 150 {
		t.Fatalf("expected progress 150 for CHAR, got %d, ok=%v", amt, ok)
	}
	if amt, ok := token.CampaignProgress("WATER"); !ok || amt != 200 {
		t.Fatalf("expected progress 200 for WATER, got %d, ok=%v", amt, ok)
	}

	// Query for non-existent campaign.
	if amt, ok := token.CampaignProgress("UNKNOWN"); ok || amt != 0 {
		t.Fatalf("expected no progress for UNKNOWN, got %d, ok=%v", amt, ok)
	}

	// Inspect campaign data and ensure it is a copy.
	camp, ok := token.Campaign("CHAR")
	if !ok {
		t.Fatalf("campaign CHAR not found")
	}
	if camp.Symbol != "CHAR" || camp.Purpose != "school" || camp.Raised != 150 {
		t.Fatalf("unexpected campaign metadata: %#v", camp)
	}
	if len(camp.Donations) != 2 || camp.Donations["alice"] != 100 || camp.Donations["bob"] != 50 {
		t.Fatalf("unexpected donations: %#v", camp.Donations)
	}

	// Modify the returned copy and ensure original data is unchanged.
	camp.Raised = 999
	camp.Donations["alice"] = 999
	camp2, _ := token.Campaign("CHAR")
	if camp2.Raised != 150 || camp2.Donations["alice"] != 100 {
		t.Fatalf("registry mutated after modifying copy: %#v", camp2)
	}

	// Ensure campaign purpose remains the first provided value.
	token.Donate("CHAR", "dave", 30, "new purpose")
	camp3, _ := token.Campaign("CHAR")
	if camp3.Purpose != "school" {
		t.Fatalf("campaign purpose changed unexpectedly: %s", camp3.Purpose)
	}
}

// TestCampaignNotFound ensures Campaign returns nil and false for unknown symbols.
func TestCampaignNotFound(t *testing.T) {
	token := NewSYN4200Token()
	if c, ok := token.Campaign("missing"); c != nil || ok {
		t.Fatalf("expected no campaign, got %#v, ok=%v", c, ok)
	}
}
