package core

import (
	"testing"
	"time"
)

func TestTradeFinanceToken(t *testing.T) {
	token := NewTradeFinanceToken()
	issue := time.Unix(0, 0)
	due := issue.Add(24 * time.Hour)
	token.RegisterDocument("doc1", "issuer", "recip", 1000, issue, due, "test")

	if err := token.FinanceDocument("doc1", "financier"); err != nil {
		t.Fatalf("finance: %v", err)
	}

	if _, ok := token.GetDocument("doc1"); !ok {
		t.Fatalf("document not found")
	}

	docs := token.ListDocuments()
	if len(docs) != 1 {
		t.Fatalf("expected 1 document, got %d", len(docs))
	}

	token.AddLiquidity("alice", 500)
	if err := token.RemoveLiquidity("alice", 200); err != nil {
		t.Fatalf("remove liquidity: %v", err)
	}
	if token.Liquidity["alice"] != 300 {
		t.Fatalf("unexpected liquidity balance: %d", token.Liquidity["alice"])
	}
}
