package optimizationnodes

import (
	"reflect"
	"sync"
	"testing"
)

// helper to create transactions quickly
func tx(hash string, fee uint64, size int) Transaction {
	return Transaction{Hash: hash, Fee: fee, Size: size}
}

func TestFeeOptimizerSortsByDensity(t *testing.T) {
	opt := &FeeOptimizer{}
	original := []Transaction{
		tx("a", 100, 100), // density 1
		tx("b", 50, 10),   // density 5
		tx("c", 90, 30),   // density 3
	}
	txs := append([]Transaction(nil), original...)

	got := opt.Optimize(txs)
	want := []Transaction{
		tx("b", 50, 10),
		tx("c", 90, 30),
		tx("a", 100, 100),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected order: %+v", got)
	}
	if !reflect.DeepEqual(txs, original) {
		t.Fatalf("input slice modified: %+v", txs)
	}
}

func TestFeeOptimizerStableSort(t *testing.T) {
	opt := &FeeOptimizer{}
	// All have same fee density (5)
	txs := []Transaction{
		tx("a", 10, 2),
		tx("b", 20, 4),
		tx("c", 30, 6),
	}
	got := opt.Optimize(txs)
	if !reflect.DeepEqual(got, txs) {
		t.Fatalf("expected stable order, got %+v", got)
	}
}

func TestFeeOptimizerZeroSize(t *testing.T) {
	opt := &FeeOptimizer{}
	txs := []Transaction{
		tx("a", 100, 0), // treated as size 1 -> density 100
		tx("b", 50, 10), // density 5
	}
	got := opt.Optimize(txs)
	want := []Transaction{
		tx("a", 100, 0),
		tx("b", 50, 10),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected order: %+v", got)
	}
}

func TestFeeOptimizerEmpty(t *testing.T) {
	opt := &FeeOptimizer{}
	if res := opt.Optimize(nil); res != nil && len(res) != 0 {
		t.Fatalf("expected nil or empty slice, got %#v", res)
	}
	if res := opt.Optimize([]Transaction{}); len(res) != 0 {
		t.Fatalf("expected empty slice, got %#v", res)
	}
}

func TestFeeOptimizerConcurrent(t *testing.T) {
	opt := &FeeOptimizer{}
	txs := []Transaction{
		tx("a", 1, 1),
		tx("b", 2, 1),
	}
	want := []Transaction{
		tx("b", 2, 1),
		tx("a", 1, 1),
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got := opt.Optimize(txs)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected order: %+v", got)
			}
		}()
	}
	wg.Wait()
}
