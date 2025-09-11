package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	sa "security_assessments"
)

type result struct {
	Name            string
	Vulnerabilities int
	Mitigations     int
	Status          string
}

func main() {
	var results []result
	walkErr := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" || strings.HasSuffix(path, "security_assesment_results.md") {
			return nil
		}
		res, aErr := sa.AnalyzeAssessment(path)
		status := "Pass"
		if aErr != nil {
			log.Printf("%s: %v", path, aErr)
			status = "Fail"
		} else if uErr := sa.UpdateTestingSection(path, res); uErr != nil {
			log.Printf("%s: %v", path, uErr)
		}
		results = append(results, result{
			Name:            strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
			Vulnerabilities: res.Vulnerabilities,
			Mitigations:     res.Mitigations,
			Status:          status,
		})
		return nil
	})
	if walkErr != nil {
		log.Fatalf("walk failed: %v", walkErr)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Name < results[j].Name })
	var b strings.Builder
	b.WriteString("| Assessment | Vulnerabilities | Mitigations | Status |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	for _, r := range results {
		b.WriteString(
			fmt.Sprintf("| %s | %d | %d | %s |\n", r.Name, r.Vulnerabilities, r.Mitigations, r.Status),
		)
	}
	if err := os.WriteFile("security_assesment_results.md", []byte(b.String()), 0o644); err != nil {
		log.Fatalf("write results: %v", err)
	}
}
