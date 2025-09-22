package tokens

import (
	"testing"
	"time"
)

func TestInvestorRegistryLifecycle(t *testing.T) {
	reg := NewInvestorRegistry()
	expiry := time.Now().Add(24 * time.Hour)
	tok, err := reg.Issue("InfraFund", "alice", 1_000, expiry)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}

	if err := reg.UpdateMetadata(tok.ID, map[string]string{"region": "APAC"}, 0.2); err != nil {
		t.Fatalf("update metadata: %v", err)
	}
	if err := reg.RecordReturn(tok.ID, 100); err != nil {
		t.Fatalf("record return: %v", err)
	}

	future := expiry.Add(48 * time.Hour)
	if err := reg.ExtendExpiry(tok.ID, future); err != nil {
		t.Fatalf("extend expiry: %v", err)
	}

	if err := reg.Deactivate(tok.ID); err != nil {
		t.Fatalf("deactivate: %v", err)
	}
	if err := reg.Activate(tok.ID); err != nil {
		t.Fatalf("activate: %v", err)
	}

	filtered := reg.FilterByAsset("InfraFund")
	if len(filtered) != 1 || filtered[0].RiskScore != 0.2 {
		t.Fatalf("unexpected filter result: %+v", filtered)
	}

	if _, ok := reg.Get(tok.ID); !ok {
		t.Fatalf("expected token to exist")
	}
}
