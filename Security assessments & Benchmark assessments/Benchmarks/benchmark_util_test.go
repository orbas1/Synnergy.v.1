package benchmarks

import (
	"strings"
	"testing"
)

func TestParseBenchmarks(t *testing.T) {
	data := []byte("BenchmarkFast-1 1 100 ns/op 0 B/op 0 allocs/op\n" +
		"BenchmarkModerate-1 1 8000 ns/op 0 B/op 0 allocs/op\n" +
		"BenchmarkSlow-1 1 2000000 ns/op 0 B/op 0 allocs/op\n")
	res, err := parseBenchmarks(data)
	if err != nil {
		t.Fatalf("parseBenchmarks returned error: %v", err)
	}
	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}
	if res[0].rating != "fast" {
		t.Errorf("expected fast rating, got %s", res[0].rating)
	}
	if res[1].rating != "moderate" {
		t.Errorf("expected moderate rating, got %s", res[1].rating)
	}
	if res[2].rating != "slow" {
		t.Errorf("expected slow rating, got %s", res[2].rating)
	}
}

func TestRepairAdvice(t *testing.T) {
	r := benchResult{name: "x", bytesPerOp: 2048, allocsPerOp: 20}
	advice := repairAdvice(r)
	if !strings.Contains(advice, "reduce allocations") || !strings.Contains(advice, "trim memory") {
		t.Errorf("unexpected advice: %s", advice)
	}
}
