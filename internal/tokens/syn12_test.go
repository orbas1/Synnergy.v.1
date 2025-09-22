package tokens

import (
	"testing"
	"time"
)

func TestSYN12Lifecycle(t *testing.T) {
	issue := time.Now().Add(-48 * time.Hour)
	maturity := time.Now().Add(24 * time.Hour)
	meta := SYN12Metadata{BillID: "TB1", Issuer: "Treasury", IssueDate: issue, Maturity: maturity, Discount: 0.95, FaceValue: 1000}
	tok := NewSYN12Token(12, "T-Bill", "TB", meta, 2)

	if err := tok.ScheduleCoupon(issue.Add(12*time.Hour), 10); err != nil {
		t.Fatalf("schedule coupon: %v", err)
	}
	if err := tok.ScheduleCoupon(issue.Add(36*time.Hour), 10); err != nil {
		t.Fatalf("schedule coupon: %v", err)
	}

	if err := tok.MarkCouponPaid(issue.Add(12*time.Hour), time.Now()); err != nil {
		t.Fatalf("mark paid: %v", err)
	}
	pending := tok.PendingCoupons()
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending coupon, got %d", len(pending))
	}

	tok.Mint("investor", 100)
	if _, err := tok.Redeem("investor", 10, time.Now()); err != ErrTreasuryNotMature {
		t.Fatalf("expected ErrTreasuryNotMature, got %v", err)
	}

	rec, err := tok.Redeem("investor", 10, maturity.Add(time.Hour))
	if err != nil {
		t.Fatalf("redeem: %v", err)
	}
	if rec.Amount != 10 || rec.Holder != "investor" {
		t.Fatalf("unexpected redemption record: %+v", rec)
	}

	history := tok.RedemptionHistory()
	if len(history) != 1 {
		t.Fatalf("expected 1 redemption, got %d", len(history))
	}
}
