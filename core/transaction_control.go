package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

// ScheduledTransaction represents a transaction scheduled for future execution.
// It can be cancelled prior to the ExecuteAt time.
type ScheduledTransaction struct {
	Tx        *Transaction
	ExecuteAt time.Time
	Canceled  bool
}

// ScheduleTransaction creates a scheduled transaction wrapper.
func ScheduleTransaction(tx *Transaction, exec time.Time) *ScheduledTransaction {
	return &ScheduledTransaction{Tx: tx, ExecuteAt: exec}
}

// CancelTransaction marks the scheduled transaction as cancelled if it
// has not yet reached its execution time.
func CancelTransaction(st *ScheduledTransaction) bool {
	if time.Now().Before(st.ExecuteAt) && !st.Canceled {
		st.Canceled = true
		return true
	}
	return false
}

// ReverseTransaction reverses a previously applied transaction. It returns an
// error if the recipient lacks sufficient funds.
func ReverseTransaction(l *Ledger, tx *Transaction) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	total, err := tx.TotalCost()
	if err != nil {
		return err
	}
	if l.balances[tx.To] < tx.Amount {
		return errors.New("insufficient recipient funds")
	}
	l.balances[tx.To] -= tx.Amount
	l.balances[tx.From] += total
	return nil
}

// ReversalRequest tracks an authority-mediated transaction reversal.
type ReversalRequest struct {
	Tx          *Transaction
	RequestedAt time.Time
	Fee         uint64
	votes       map[string]bool
}

// RequestReversal freezes the recipient's funds and records a reversal request.
// The recipient must cover the amount plus return gas fee.
func RequestReversal(l *Ledger, tx *Transaction, fee uint64) (*ReversalRequest, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if tx.Amount > math.MaxUint64-fee {
		return nil, ErrTransactionOverflow
	}
	total := tx.Amount + fee
	if l.balances[tx.To] < total {
		return nil, errors.New("insufficient funds to freeze for reversal")
	}
	l.balances[tx.To] -= total
	l.frozen[tx.To] += total
	return &ReversalRequest{Tx: tx, RequestedAt: time.Now(), Fee: fee, votes: make(map[string]bool)}, nil
}

// Vote records an authority node's decision on the reversal request.
func (r *ReversalRequest) Vote(authorityID string, approve bool) {
	r.votes[authorityID] = approve
}

const reversalWindow = 30 * 24 * time.Hour

// FinalizeReversal executes the compensating transaction if enough authority
// nodes approve within the time window. Frozen funds are used to pay the
// reversal and associated fee.
func FinalizeReversal(l *Ledger, r *ReversalRequest, required int) error {
	approvals := 0
	for _, v := range r.votes {
		if v {
			approvals++
		}
	}
	if approvals < required {
		return errors.New("insufficient authority approvals")
	}
	if time.Since(r.RequestedAt) > reversalWindow {
		RejectReversal(l, r)
		return errors.New("reversal request expired")
	}
	if r.Tx.Amount > math.MaxUint64-r.Fee {
		return ErrTransactionOverflow
	}
	total := r.Tx.Amount + r.Fee
	l.mu.Lock()
	l.balances[r.Tx.To] += total
	l.frozen[r.Tx.To] -= total
	l.mu.Unlock()
	revTx := NewTransaction(r.Tx.To, r.Tx.From, r.Tx.Amount, r.Fee, 0)
	return l.ApplyTransaction(revTx)
}

// RejectReversal releases frozen funds when a reversal request fails.
func RejectReversal(l *Ledger, r *ReversalRequest) {
	if r.Tx.Amount > math.MaxUint64-r.Fee {
		return
	}
	total := r.Tx.Amount + r.Fee
	l.mu.Lock()
	l.balances[r.Tx.To] += total
	l.frozen[r.Tx.To] -= total
	l.mu.Unlock()
}

// ConvertToPrivate encrypts the transaction using AES-GCM with the provided key.
func ConvertToPrivate(tx *Transaction, key []byte) (*PrivateTransaction, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	payload := gcm.Seal(nil, nonce, b, nil)
	return &PrivateTransaction{Payload: payload, Nonce: nonce}, nil
}

// Decrypt decrypts the private transaction using the same key used for
// encryption.
func (pt *PrivateTransaction) Decrypt(key []byte) (*Transaction, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	b, err := gcm.Open(nil, pt.Nonce, pt.Payload, nil)
	if err != nil {
		return nil, err
	}
	var tx Transaction
	if err := json.Unmarshal(b, &tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

// Receipt captures the outcome of a transaction.
type Receipt struct {
	TxID      string
	Status    string
	Timestamp int64
	Details   string
}

// GenerateReceipt creates a receipt for the given transaction ID and status.
func GenerateReceipt(txID, status, details string) Receipt {
	return Receipt{TxID: txID, Status: status, Timestamp: time.Now().Unix(), Details: details}
}

// ReceiptStore provides thread-safe storage for receipts.
type ReceiptStore struct {
	mu        sync.RWMutex
	data      map[string]Receipt
	retention time.Duration
	clock     func() time.Time
}

// ReceiptStoreOption customises receipt store behaviour.
type ReceiptStoreOption func(*ReceiptStore)

// WithReceiptRetention limits receipt lifetime to the provided duration. A zero
// duration preserves all receipts.
func WithReceiptRetention(d time.Duration) ReceiptStoreOption {
	return func(rs *ReceiptStore) {
		rs.retention = d
	}
}

// WithReceiptClock sets the clock used for retention decisions, simplifying
// deterministic testing.
func WithReceiptClock(clock func() time.Time) ReceiptStoreOption {
	return func(rs *ReceiptStore) {
		if clock != nil {
			rs.clock = clock
		}
	}
}

// NewReceiptStore constructs an empty receipt store.
func NewReceiptStore(opts ...ReceiptStoreOption) *ReceiptStore {
	rs := &ReceiptStore{
		data:  make(map[string]Receipt),
		clock: time.Now,
	}
	for _, opt := range opts {
		opt(rs)
	}
	return rs
}

// Store saves a receipt in the store.
func (rs *ReceiptStore) Store(r Receipt) {
	rs.mu.Lock()
	rs.data[r.TxID] = r
	rs.purgeLocked()
	rs.mu.Unlock()
}

// Get retrieves a receipt by transaction ID.
func (rs *ReceiptStore) Get(id string) (Receipt, bool) {
	rs.mu.RLock()
	r, ok := rs.data[id]
	rs.mu.RUnlock()
	return r, ok
}

// List returns all receipts ordered by timestamp then TxID for deterministic
// CLI and API responses.
func (rs *ReceiptStore) List() []Receipt {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	res := make([]Receipt, 0, len(rs.data))
	for _, r := range rs.data {
		res = append(res, r)
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].Timestamp == res[j].Timestamp {
			return res[i].TxID < res[j].TxID
		}
		return res[i].Timestamp < res[j].Timestamp
	})
	return res
}

// Search returns all receipts containing the keyword in any field.
func (rs *ReceiptStore) Search(keyword string) []Receipt {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	var res []Receipt
	for _, r := range rs.data {
		if strings.Contains(r.TxID, keyword) || strings.Contains(r.Status, keyword) || strings.Contains(r.Details, keyword) {
			res = append(res, r)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].Timestamp == res[j].Timestamp {
			return res[i].TxID < res[j].TxID
		}
		return res[i].Timestamp < res[j].Timestamp
	})
	return res
}

// PurgeExpired removes receipts older than the configured retention window and
// returns the count of records removed.
func (rs *ReceiptStore) PurgeExpired() int {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.purgeLocked()
}

func (rs *ReceiptStore) purgeLocked() int {
	if rs.retention <= 0 {
		return 0
	}
	now := rs.clock()
	removed := 0
	for id, r := range rs.data {
		if now.Sub(time.Unix(r.Timestamp, 0)) > rs.retention {
			delete(rs.data, id)
			removed++
		}
	}
	return removed
}
