package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Transaction represents a transfer of Synthron between accounts.
//
// Transactions are signed payloads that move coins from one address to
// another.  They include a fee and timestamp so that they can be ordered
// deterministically by the consensus engine.
type Transaction struct {
	ID        string
	From      string
	To        string
	Amount    uint64
	Fee       uint64
	Nonce     uint64
	Timestamp int64
	Signature []byte
}

// NewTransaction creates a new unsigned transaction with the provided
// parameters.  The ID is derived from a hash of the core fields and can be
// reproduced deterministically prior to signing.
func NewTransaction(from, to string, amount, fee, nonce uint64) *Transaction {
	tx := &Transaction{From: from, To: to, Amount: amount, Fee: fee, Nonce: nonce, Timestamp: time.Now().Unix()}
	tx.ID = tx.Hash()
	return tx
}

// Hash returns the hex-encoded hash of the transaction contents excluding the
// signature.  It is used as the message for signing and verification.
func (t *Transaction) Hash() string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s%s%d%d%d%d", t.From, t.To, t.Amount, t.Fee, t.Nonce, t.Timestamp)))
	return hex.EncodeToString(h[:])
}
