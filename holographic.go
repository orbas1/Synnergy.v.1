package synnergy

// HolographicFrame represents holographically encoded data.
type HolographicFrame struct {
	ID     string
	Shards [][]byte
}

// SplitHolographic divides data into n shards for redundant distribution.
func SplitHolographic(id string, data []byte, n int) HolographicFrame {
	if n <= 0 {
		return HolographicFrame{ID: id}
	}
	shards := make([][]byte, n)
	size := (len(data) + n - 1) / n
	for i := 0; i < n; i++ {
		start := i * size
		end := start + size
		if end > len(data) {
			end = len(data)
		}
		shard := make([]byte, end-start)
		copy(shard, data[start:end])
		shards[i] = shard
	}
	return HolographicFrame{ID: id, Shards: shards}
}

// ReconstructHolographic recombines shards into the original byte slice.
func ReconstructHolographic(frame HolographicFrame) []byte {
	total := 0
	for _, s := range frame.Shards {
		total += len(s)
	}
	out := make([]byte, 0, total)
	for _, s := range frame.Shards {
		out = append(out, s...)
	}
	return out
}
