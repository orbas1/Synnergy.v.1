package benchmarks

var benchmarkSink int

func performComputation(iterations int) int {
	sum := 0
	for i := 0; i < iterations; i++ {
		sum += i
	}
	return sum
}
