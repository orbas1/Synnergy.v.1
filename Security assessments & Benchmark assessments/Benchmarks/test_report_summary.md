# Test Run Summary

A full test sweep was attempted with extended timeouts.

## Results

- **synnergy/core**: all tests passed (`go test ./core -timeout 30m`).
- **synnergy/cli**: full suite ran for ~189s before manual interruption; individual JSON command tests like `TestDAOTokenMintJSON` and `TestElectedNodeCreateJSON` complete in ~0.05s.
- `benchreport` analyzed 14 functions across the benchmarking module.
- `runbench` regenerated `Benchmarks_full_report_and_assessment.md`.

The CLI package houses numerous command tests that significantly prolong the overall test runtime.
