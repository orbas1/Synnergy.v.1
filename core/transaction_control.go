package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
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
	total := tx.Amount + tx.Fee
	if l.balances[tx.To] < tx.Amount {
		return errors.New("insufficient recipient funds")
	}
	l.balances[tx.To] -= tx.Amount
	l.balances[tx.From] += total
	return nil
}

// ConvertToPrivate encrypts the transaction using AES-CTR with the provided key.
func ConvertToPrivate(tx *Transaction, key []byte) (*PrivateTransaction, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, nonce)
	payload := make([]byte, len(b))
	stream.XORKeyStream(payload, b)
	return &PrivateTransaction{Payload: payload, Nonce: nonce}, nil
}

// Decrypt decrypts the private transaction using the same key used for
// encryption.
func (pt *PrivateTransaction) Decrypt(key []byte) (*Transaction, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, pt.Nonce)
	b := make([]byte, len(pt.Payload))
	stream.XORKeyStream(b, pt.Payload)
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
	mu   sync.RWMutex
	data map[string]Receipt
}

// NewReceiptStore constructs an empty receipt store.
func NewReceiptStore() *ReceiptStore {
	return &ReceiptStore{data: make(map[string]Receipt)}
}

// Store saves a receipt in the store.
func (rs *ReceiptStore) Store(r Receipt) {
	rs.mu.Lock()
	rs.data[r.TxID] = r
	rs.mu.Unlock()
}

// Get retrieves a receipt by transaction ID.
func (rs *ReceiptStore) Get(id string) (Receipt, bool) {
	rs.mu.RLock()
	r, ok := rs.data[id]
	rs.mu.RUnlock()
	return r, ok
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
	return res
}
