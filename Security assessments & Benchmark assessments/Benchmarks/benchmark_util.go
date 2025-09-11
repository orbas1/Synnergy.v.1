package benchmarks

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

var benchmarkSink int

func performComputation(iterations int) int {
	sum := 0
	for i := 0; i < iterations; i++ {
		sum += i
	}
	return sum
}

// benchLine matches the standard Go benchmark output produced with -benchmem.
// Example line:
// BenchmarkFoo-8        1   123 ns/op   45 B/op   2 allocs/op
var benchLine = regexp.MustCompile(`^(Benchmark\S+)\s+\d+\s+(\d+(?:\.\d+)?)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op`)

type benchResult struct {
	name        string
	nsPerOp     float64
	bytesPerOp  int
	allocsPerOp int
	opsPerSec   float64
	rating      string
}

// repairAdvice returns a human readable optimisation tip based on
// allocation and memory characteristics of the benchmark result.
func repairAdvice(r benchResult) string {
	advice := "profile to optimise algorithmic efficiency"
	if r.allocsPerOp > 10 {
		advice = "reduce allocations and reuse objects"
	}
	if r.bytesPerOp > 1024 {
		if advice != "profile to optimise algorithmic efficiency" {
			advice += "; "
		} else {
			advice = ""
		}
		advice += "trim memory usage"
	}
	if advice == "" {
		advice = "profile to optimise algorithmic efficiency"
	}
	return advice
}

// parseBenchmarks converts raw `go test` benchmark output into structured
// results enriched with throughput and a qualitative rating. Ratings are
// assigned using simple thresholds on operations per second: >500k "fast",
// >100k "moderate", otherwise "slow".
func parseBenchmarks(out []byte) ([]benchResult, error) {
	scanner := bufio.NewScanner(bytes.NewReader(out))
	var results []benchResult
	for scanner.Scan() {
		matches := benchLine.FindStringSubmatch(scanner.Text())
		if len(matches) != 5 {
			continue
		}
		ns, err1 := strconv.ParseFloat(matches[2], 64)
		b, err2 := strconv.Atoi(matches[3])
		a, err3 := strconv.Atoi(matches[4])
		if err1 != nil || err2 != nil || err3 != nil {
			continue
		}
		ops := 1e9 / ns
		rating := "slow"
		if ops > 500000 {
			rating = "fast"
		} else if ops > 100000 {
			rating = "moderate"
		}
		results = append(results, benchResult{matches[1], ns, b, a, ops, rating})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no benchmarks found")
	}
	return results, nil
}

// GenerateBenchmarkReport runs benchmarks across the repository and appends a
// markdown table of results along with a brief analysis to the provided output file.
func GenerateBenchmarkReport(outputPath string) error {
	absOut, err := filepath.Abs(outputPath)
	if err != nil {
		return err
	}
	benchDir := filepath.Dir(absOut)
	moduleDir := filepath.Dir(benchDir)
	rootDir := filepath.Dir(moduleDir)
	cmd := exec.Command("go", "test", "-run=^$", "-bench=.", "-benchtime=1x", "-benchmem", "./...")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("running benchmarks: %w\n%s", err, output)
	}

	results, err := parseBenchmarks(output)
	if err != nil {
		return err
	}

	sort.Slice(results, func(i, j int) bool { return results[i].nsPerOp < results[j].nsPerOp })
	fastest := results[0]
	slowest := results[len(results)-1]

	var buf bytes.Buffer
	buf.WriteString("| Benchmark | ns/op | ops/s | B/op | allocs/op | Rating |\n")
	buf.WriteString("|-----------|------|-------|------|-----------|--------|\n")
	var totalOps float64
	var slowList []benchResult
	for _, r := range results {
		totalOps += r.opsPerSec
		fmt.Fprintf(&buf, "| %s | %.2f | %.0f | %d | %d | %s |\n", r.name, r.nsPerOp, r.opsPerSec, r.bytesPerOp, r.allocsPerOp, r.rating)
		if r.rating == "slow" {
			slowList = append(slowList, r)
		}
	}
	avgOps := totalOps / float64(len(results))
	overall := "slow"
	if avgOps > 500000 {
		overall = "fast"
	} else if avgOps > 100000 {
		overall = "moderate"
	}
	analysis := fmt.Sprintf("\n### Analysis\n\nFastest benchmark **%s** at %.2f ns/op (%.0f ops/s). "+
		"Slowest benchmark **%s** at %.2f ns/op (%.0f ops/s).\n"+
		"Average throughput %.0f ops/s, indicating %s overall performance. "+
		"The slowest benchmark suggests a ceiling of %.2f ns/op (~%.0f ops/s).\n"+
		"The ns/op metric reflects the time each operation takes; lower ns/op and higher ops/s indicate better performance.\n",
		fastest.name, fastest.nsPerOp, fastest.opsPerSec,
		slowest.name, slowest.nsPerOp, slowest.opsPerSec,
		avgOps, overall, slowest.nsPerOp, slowest.opsPerSec)

	upgrade := "\n### Upgrade Plan\n\n"
	if len(slowList) > 0 {
		upgrade += "The following benchmarks are classified as slow and may need optimisation:\n\n"
		for _, s := range slowList {
			upgrade += fmt.Sprintf("- %s (%.2f ns/op)\n", s.name, s.nsPerOp)
		}
	} else {
		upgrade += "All benchmarks meet performance targets; no immediate upgrades required.\n"
	}

	bottleneck := "\n### Bottleneck Analysis and Repair Plan\n\n"
	if len(slowList) > 0 {
		for _, s := range slowList {
			bottleneck += fmt.Sprintf("- %s: %s\n", s.name, repairAdvice(s))
		}
	} else {
		bottleneck += "No bottlenecks detected.\n"
	}

	intro := "# Benchmark Full Report and Assessment\n\n" +
		"This document provides an enterprise-grade view of benchmark performance across the Synnergy codebase. " +
		"Benchmarks are executed with `go test -bench . -benchmem ./...`.\n\n" +
		"Run `cd 'Security assessments & Benchmark assessments' && go run ./Benchmarks/cmd/runbench` " +
		"to refresh results. The command rewrites this file with the latest measurements.\n\n" +
		"The `ns/op` column reports average time per operation, while `ops/s` derives throughput. " +
		"Lower `ns/op` and higher `ops/s` indicate better performance.\n\n"

	f, err := os.OpenFile(absOut, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(intro + "## Automated Benchmark Results\n\n"); err != nil {
		return err
	}
	if _, err := f.WriteString(buf.String()); err != nil {
		return err
	}
	_, err = f.WriteString(analysis + upgrade + bottleneck)
	return err
}
