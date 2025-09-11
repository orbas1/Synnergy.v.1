# Benchmark Results

The following benchmarks were run with:

```
go test -bench . ./...
```

## Test Summary

All unit tests passed when executing:

```
go test ./...
```

## synnergy package

| Benchmark | Iterations | ns/op |
| --- | --- | --- |
| BenchmarkTransactionManagerLockAndMint | 367810 | 2873 |
| BenchmarkTransactionManagerBurnAndRelease | 539388 | 2880 |
| BenchmarkTransactionManagerListTransactions | 10000 | 599597 |
| BenchmarkTransactionManagerGetTransaction | 3155931 | 380.1 |

## synnergy/core package

| Benchmark | Iterations | ns/op |
| --- | --- | --- |
| BenchmarkTransactionHash | 1906989 | 627.8 |
| BenchmarkLedgerApplyTransaction | 2172804 | 518.6 |
| BenchmarkNFTMarketplaceMint | 1000000 | 2074 |

## synnergy/internal/tokens package

| Benchmark | Iterations | ns/op |
| --- | --- | --- |
| BenchmarkBaseTokenTransfer | 30524721 | 39.35 |
| BenchmarkRegistryInfo | 11217840 | 107.9 |

## security_assessments_benchmarks/Benchmarks package

| Benchmark | Iterations | ns/op |
| --- | --- | --- |
| BenchmarkAllTokenStandardBenchmarks | 1 | 1661 |
| BenchmarkAuthorityNodeBenchmarks | 1 | 826 |
| BenchmarkFullReportAndAssessment | 1 | 783 |
| BenchmarkCharityBenchmarks | 1 | 973 |
| BenchmarkCoinBenchmarks | 1 | 782 |
| BenchmarkComplianceBenchmarks | 1 | 838 |
| BenchmarkConsensusBenchmarks | 1 | 820 |
| BenchmarkContractBenchmarks | 1 | 812 |
| BenchmarkGovernanceBenchmarks | 1 | 767 |
| BenchmarkHighAvailabilityBenchmarks | 1 | 819 |
| BenchmarkLedgerBenchmarks | 1 | 10378 |
| BenchmarkLoanpoolBenchmarks | 1 | 831 |
| BenchmarkNetworkBenchmarks | 1 | 832 |
| BenchmarkNodeBenchmarks | 1 | 759 |
| BenchmarkOpcodeBenchmarks | 1 | 1234 |
| BenchmarkSecurityBenchmarks | 1 | 1465 |
| BenchmarkSpeedBenchmarks | 1 | 1346 |
| BenchmarkStorageBenchmarks | 1 | 795 |
| BenchmarkVmBenchmarks | 1 | 1320 |
| BenchmarkValidationBenchmarks | 1 | 881 |
| BenchmarkWalletBenchmarks | 1 | 792 |
| BenchmarkAiBenchmarks | 1 | 801 |
| BenchmarkTransactionsBenchmarks | 1 | 736 |

Benchmarks executed on `Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz`.

## Slowest Benchmarks

The `Security assessments & Benchmark assessments/Benchmarks/cmd/benchreport` utility enumerates all functions and runs benchmarks package by package:

```
go run "Security assessments & Benchmark assessments/Benchmarks/cmd/benchreport"
```

The slowest benchmarks from the latest run are:

| Package | Benchmark | ns/op | B/op | allocs/op |
| --- | --- | --- | --- | --- |
| synnergy | BenchmarkTransactionManagerListTransactions | 611597 | 1286156 | 1 |
| synnergy | BenchmarkTransactionManagerLockAndMint | 2824 | 585 | 7 |
| synnergy | BenchmarkTransactionManagerBurnAndRelease | 2763 | 822 | 9 |
| synnergy/core | BenchmarkNFTMarketplaceMint | 1869 | 431 | 6 |
| synnergy/core | BenchmarkTransactionHash | 677 | 176 | 5 |
| synnergy | BenchmarkTransactionManagerGetTransaction | 538 | 0 | 0 |
| synnergy/core | BenchmarkLedgerApplyTransaction | 534 | 96 | 8 |
| synnergy/internal/tokens | BenchmarkRegistryInfo | 107 | 0 | 0 |
| synnergy/internal/tokens | BenchmarkBaseTokenTransfer | 39 | 0 | 0 |

