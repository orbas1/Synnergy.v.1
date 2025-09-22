package tokens

import (
	"errors"
	"sync"
	"time"
)

// SYN12Metadata defines treasury bill properties for CBDC instruments.
type SYN12Metadata struct {
	BillID    string
	Issuer    string
	IssueDate time.Time
	Maturity  time.Time
	Discount  float64
	FaceValue uint64
}

// SYN12Token represents a tokenised treasury bill.
type SYN12Token struct {
	*BaseToken
	mu          sync.RWMutex
	Metadata    SYN12Metadata
	coupons     []SYN12Coupon
	redemptions []RedemptionRecord
}

// NewSYN12Token creates a new SYN12 token with the provided metadata.
func NewSYN12Token(id TokenID, name, symbol string, meta SYN12Metadata, decimals uint8) *SYN12Token {
	return &SYN12Token{
		BaseToken: NewBaseToken(id, name, symbol, decimals),
		Metadata:  meta,
	}
}

// SYN12Coupon represents a scheduled coupon payment.
type SYN12Coupon struct {
	DueDate time.Time
	Amount  uint64
	Paid    bool
	PaidAt  time.Time
}

// RedemptionRecord captures redemption operations.
type RedemptionRecord struct {
	Holder      string
	Amount      uint64
	RequestedAt time.Time
	ProcessedAt time.Time
}

var (
	// ErrCouponOutOfRange indicates an invalid coupon schedule.
	ErrCouponOutOfRange = errors.New("coupon due date outside treasury term")
	// ErrCouponAlreadyRecorded indicates a duplicate coupon.
	ErrCouponAlreadyRecorded = errors.New("coupon already scheduled")
	// ErrCouponNotFound is returned when updating a non-existent coupon.
	ErrCouponNotFound = errors.New("coupon not found")
	// ErrTreasuryNotMature indicates redemption attempted before maturity.
	ErrTreasuryNotMature = errors.New("treasury bill has not matured")
)

// ScheduleCoupon adds a coupon payment to the treasury bill lifecycle.
func (t *SYN12Token) ScheduleCoupon(due time.Time, amount uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if due.Before(t.Metadata.IssueDate) || due.After(t.Metadata.Maturity) {
		return ErrCouponOutOfRange
	}
	for _, c := range t.coupons {
		if c.DueDate.Equal(due) {
			return ErrCouponAlreadyRecorded
		}
	}
	t.coupons = append(t.coupons, SYN12Coupon{DueDate: due, Amount: amount})
	return nil
}

// MarkCouponPaid flags the coupon as satisfied.
func (t *SYN12Token) MarkCouponPaid(due time.Time, paidAt time.Time) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, c := range t.coupons {
		if c.DueDate.Equal(due) {
			t.coupons[i].Paid = true
			t.coupons[i].PaidAt = paidAt
			return nil
		}
	}
	return ErrCouponNotFound
}

// PendingCoupons returns unpaid coupon records.
func (t *SYN12Token) PendingCoupons() []SYN12Coupon {
	t.mu.RLock()
	defer t.mu.RUnlock()
	pending := make([]SYN12Coupon, 0)
	for _, c := range t.coupons {
		if !c.Paid {
			pending = append(pending, c)
		}
	}
	return pending
}

// Redeem processes redemptions only after maturity.
func (t *SYN12Token) Redeem(holder string, amount uint64, now time.Time) (RedemptionRecord, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if now.Before(t.Metadata.Maturity) {
		return RedemptionRecord{}, ErrTreasuryNotMature
	}
	if err := t.BaseToken.Burn(holder, amount); err != nil {
		return RedemptionRecord{}, err
	}
	rec := RedemptionRecord{Holder: holder, Amount: amount, RequestedAt: now, ProcessedAt: time.Now()}
	t.redemptions = append(t.redemptions, rec)
	return rec, nil
}

// RedemptionHistory exposes prior redemption operations.
func (t *SYN12Token) RedemptionHistory() []RedemptionRecord {
	t.mu.RLock()
	defer t.mu.RUnlock()
	hist := make([]RedemptionRecord, len(t.redemptions))
	copy(hist, t.redemptions)
	return hist
}
