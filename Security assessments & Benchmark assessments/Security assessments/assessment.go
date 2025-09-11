package security

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type AssessmentResult struct {
	Vulnerabilities int
	Mitigations     int
}

// AnalyzeAssessment reads the markdown file and ensures required sections and
// a minimum number of bullet points for vulnerabilities and mitigations.
func AnalyzeAssessment(filename string) (AssessmentResult, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return AssessmentResult{}, fmt.Errorf("failed to read %s: %w", filename, err)
	}
	content := string(data)
	required := []string{
		"## Overview",
		"## Potential Vulnerabilities",
		"## Mitigation Strategies",
		"## Security Testing",
		"## Conclusion",
	}
	for _, s := range required {
		if !strings.Contains(content, s) {
			return AssessmentResult{}, fmt.Errorf("%s missing section %q", filename, s)
		}
	}
	if strings.Contains(content, "TODO") || strings.Contains(content, "TBD") {
		return AssessmentResult{}, fmt.Errorf("%s contains unfinished placeholder text", filename)
	}
	scanner := bufio.NewScanner(strings.NewReader(content))
	var current string
	res := AssessmentResult{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "##") {
			current = line
			continue
		}
		if strings.HasPrefix(line, "- ") {
			switch current {
			case "## Potential Vulnerabilities":
				res.Vulnerabilities++
			case "## Mitigation Strategies":
				res.Mitigations++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return res, fmt.Errorf("scan failed for %s: %w", filename, err)
	}
	if res.Vulnerabilities == 0 {
		return res, fmt.Errorf("%s has no vulnerabilities listed", filename)
	}
	if res.Mitigations == 0 {
		return res, fmt.Errorf("%s has no mitigations listed", filename)
	}
	return res, nil
}

// UpdateTestingSection appends or replaces a line in the Security Testing
// section summarizing automated validation results.
func UpdateTestingSection(filename string, res AssessmentResult) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filename, err)
	}
	lines := strings.Split(string(data), "\n")
	out := make([]string, 0, len(lines)+1)
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		out = append(out, line)
		if strings.TrimSpace(line) == "## Security Testing" {
			if i+1 < len(lines) {
				out = append(out, lines[i+1])
				i++
			}
			msg := fmt.Sprintf("Automated validation on %s verified %d vulnerabilities and %d mitigations.",
				time.Now().Format("2006-01-02"), res.Vulnerabilities, res.Mitigations)
			if i+1 < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i+1]), "Automated validation") {
				i++ // skip existing automated line
			}
			out = append(out, msg)
		}
	}
	return os.WriteFile(filename, []byte(strings.Join(out, "\n")), 0o644)
}
