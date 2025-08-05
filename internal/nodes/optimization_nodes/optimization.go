package optimizationnodes

import (
	"sort"
	"sync"
)

// OptimizationNode defines a component that can reorder transactions to
// improve network performance.
type OptimizationNode interface {
	// Optimize returns a reordered copy of the supplied transactions.
	// Implementations may apply custom heuristics such as maximising fee
	// per byte or minimising latency.
	Optimize(txs []Transaction) []Transaction
}

// FeeOptimizer orders transactions by fee density (fee per byte) in descending
// order. It is safe for concurrent use.
type FeeOptimizer struct {
	mu sync.Mutex
}

// Optimize implements the OptimizationNode interface.
// Transactions are sorted by fee density with higher paying transactions first.
func (o *FeeOptimizer) Optimize(txs []Transaction) []Transaction {
	o.mu.Lock()
	defer o.mu.Unlock()

	out := make([]Transaction, len(txs))
	copy(out, txs)
	sort.SliceStable(out, func(i, j int) bool {
		fi := float64(out[i].Fee)
		fj := float64(out[j].Fee)
		si := float64(out[i].Size)
		sj := float64(out[j].Size)
		if si == 0 {
			si = 1
		}
		if sj == 0 {
			sj = 1
		}
		return fi/si > fj/sj
	})
	return out
}
