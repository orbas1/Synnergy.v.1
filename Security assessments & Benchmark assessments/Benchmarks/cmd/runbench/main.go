package main

import (
	"flag"
	"log"

	benchmarks "security_assessments_benchmarks/Benchmarks"
)

func main() {
	out := flag.String("out", "Benchmarks/Benchmarks_full_report_and_assessment.md", "output markdown file")
	flag.Parse()
	if err := benchmarks.GenerateBenchmarkReport(*out); err != nil {
		log.Fatalf("failed to generate benchmark report: %v", err)
	}
	log.Printf("benchmark report written to %s", *out)
}
