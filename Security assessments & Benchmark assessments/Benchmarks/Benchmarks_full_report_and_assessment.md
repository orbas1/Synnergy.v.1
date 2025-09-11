# Benchmark Full Report and Assessment

The benchmarks in this report measure performance across the security assessment suite. Each function was executed with `go test -run=^$ -bench . -benchtime=1x` and metrics were captured for execution time and allocations. All runs were performed on an `Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz`.

| Benchmark | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| BenchmarkAllTokenStandardBenchmarks | 2511 | 0 | 0 |
| BenchmarkAuthorityNodeBenchmarks | 1580 | 0 | 0 |
| BenchmarkFullReportAndAssessment | 1977 | 0 | 0 |
| BenchmarkCharityBenchmarks | 855 | 0 | 0 |
| BenchmarkCoinBenchmarks | 1413 | 0 | 0 |
| BenchmarkComplianceBenchmarks | 839 | 0 | 0 |
| BenchmarkConsensusBenchmarks | 955 | 0 | 0 |
| BenchmarkContractBenchmarks | 835 | 0 | 0 |
| BenchmarkGovernanceBenchmarks | 1102 | 0 | 0 |
| BenchmarkHighAvailabilityBenchmarks | 850 | 0 | 0 |
| BenchmarkLedgerBenchmarks | 1411 | 0 | 0 |
| BenchmarkLoanpoolBenchmarks | 1558 | 0 | 0 |
| BenchmarkNetworkBenchmarks | 847 | 0 | 0 |
| BenchmarkNodeBenchmarks | 756 | 0 | 0 |
| BenchmarkOpcodeBenchmarks | 910 | 0 | 0 |
| BenchmarkSecurityBenchmarks | 1384 | 0 | 0 |
| BenchmarkSpeedBenchmarks | 797 | 0 | 0 |
| BenchmarkStorageBenchmarks | 812 | 0 | 0 |
| BenchmarkVmBenchmarks | 1355 | 0 | 0 |
| BenchmarkValidationBenchmarks | 1556 | 0 | 0 |
| BenchmarkWalletBenchmarks | 766 | 0 | 0 |
| BenchmarkAiBenchmarks | 933 | 0 | 0 |
| BenchmarkTransactionsBenchmarks | 762 | 0 | 0 |

These figures provide a baseline for future optimisations. Contributors should monitor this report when evaluating the impact of changes on system performance.
