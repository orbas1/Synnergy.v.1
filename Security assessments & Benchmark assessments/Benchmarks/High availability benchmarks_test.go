package benchmarks

import "testing"

func BenchmarkHighAvailabilityBenchmarks(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSink += performComputation(1000)
	}
}
