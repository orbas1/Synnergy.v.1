package core

import (
	"sync"
	"testing"
	"time"
)

func TestTradeFinanceTokenBasic(t *testing.T) {
	token := NewTradeFinanceToken()
	issue := time.Unix(0, 0)
	due := issue.Add(24 * time.Hour)
	if err := token.RegisterDocument("doc1", "issuer", "recip", 1000, issue, due, "test"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := token.RegisterDocument("doc1", "issuer", "recip", 1000, issue, due, "test"); err != ErrDocumentExists {
		t.Fatalf("expected ErrDocumentExists, got %v", err)
	}
	if err := token.FinanceDocument("doc1", "financier"); err != nil {
		t.Fatalf("finance: %v", err)
	}
	if err := token.FinanceDocument("doc1", "financier"); err != ErrDocumentAlreadyFinanced {
		t.Fatalf("expected ErrDocumentAlreadyFinanced, got %v", err)
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
	if err := token.RemoveLiquidity("alice", 400); err != ErrInsufficientLiquidity {
		t.Fatalf("expected ErrInsufficientLiquidity, got %v", err)
	}
}

func TestTradeFinanceTokenConcurrentLiquidity(t *testing.T) {
	token := NewTradeFinanceToken()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token.AddLiquidity("addr", 1)
		}()
	}
	wg.Wait()
	if token.Liquidity["addr"] != 50 {
		t.Fatalf("expected 50 liquidity, got %d", token.Liquidity["addr"])
	}
}
