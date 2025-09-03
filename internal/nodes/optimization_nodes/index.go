package optimizationnodes

import "sort"

// Transaction represents the minimal transaction information required for
// optimisation decisions. Only fields relevant to ordering are included.
type Transaction struct {
	Hash string // unique transaction identifier
	Fee  uint64 // total fee offered by the transaction
	Size int    // approximate size in bytes
}

// Sorter provides helpers for selecting transactions to include in a block
// based on fee density (fee per byte).  This simple implementation is
// deterministic and safe for concurrent use if the input slice is not shared.
type Sorter struct{}

// Select returns a subset of transactions ordered by highest fee density while
// keeping the total size under maxBytes.
func (Sorter) Select(txns []Transaction, maxBytes int) []Transaction {
	sort.Slice(txns, func(i, j int) bool {
		f1 := float64(txns[i].Fee) / float64(txns[i].Size)
		f2 := float64(txns[j].Fee) / float64(txns[j].Size)
		return f1 > f2
	})
	var total int
	var out []Transaction
	for _, tx := range txns {
		if total+tx.Size > maxBytes {
			continue
		}
		out = append(out, tx)
		total += tx.Size
	}
	return out
}
