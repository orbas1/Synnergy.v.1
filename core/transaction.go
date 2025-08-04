package core

// Transaction represents a transfer of Synthron between accounts.
type Transaction struct {
	From      string
	To        string
	Amount    uint64
	Nonce     uint64
	Signature []byte
}

// NewTransaction creates a new unsigned transaction.
func NewTransaction(from, to string, amount, nonce uint64) *Transaction {
	return &Transaction{From: from, To: to, Amount: amount, Nonce: nonce}
}
