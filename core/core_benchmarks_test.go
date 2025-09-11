package core

import "testing"

var stringSink string

func BenchmarkTransactionHash(b *testing.B) {
	tx := &Transaction{
		From:      "alice",
		To:        "bob",
		Amount:    1,
		Fee:       1,
		Nonce:     1,
		Timestamp: 1,
		Type:      TxTypeTransfer,
	}
	var r string
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r = tx.Hash()
	}
	stringSink = r
}

func BenchmarkLedgerApplyTransaction(b *testing.B) {
	l := NewLedger()
	l.Credit("alice", uint64(b.N)+1)
	tx := NewTransaction("alice", "bob", 1, 0, 0)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := l.ApplyTransaction(tx); err != nil {
			b.Fatal(err)
		}
	}
}
