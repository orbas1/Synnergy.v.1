# Benchmark Full Report and Assessment

This document provides an enterprise-grade view of benchmark performance across the Synnergy codebase. Benchmarks are executed with `go test -bench . -benchmem ./...`.

Run `cd 'Security assessments & Benchmark assessments' && go run ./Benchmarks/cmd/runbench` to refresh results. The command rewrites this file with the latest measurements.

The `ns/op` column reports average time per operation, while `ops/s` derives throughput. Lower `ns/op` and higher `ops/s` indicate better performance.

## Automated Benchmark Results

| Benchmark | ns/op | ops/s | B/op | allocs/op | Rating |
|-----------|------|-------|------|-----------|--------|
| BenchmarkTransactionManagerGetTransaction-5 | 1108.00 | 902527 | 0 | 0 | fast |
| BenchmarkBaseTokenTransfer-5 | 2596.00 | 385208 | 0 | 0 | moderate |
| BenchmarkRegistryInfo-5 | 3781.00 | 264480 | 0 | 0 | moderate |
| BenchmarkTransactionManagerListTransactions-5 | 4619.00 | 216497 | 128 | 1 | moderate |
| BenchmarkLedgerApplyTransaction-5 | 8049.00 | 124239 | 80 | 6 | moderate |
| BenchmarkTransactionManagerBurnAndRelease-5 | 199231.00 | 5019 | 2360 | 14 | slow |
| BenchmarkTransactionManagerLockAndMint-5 | 200195.00 | 4995 | 2600 | 17 | slow |
| BenchmarkTransactionHash-5 | 1400781.00 | 714 | 1128 | 10 | slow |
| BenchmarkNFTMarketplaceMint-5 | 2217720.00 | 451 | 164920 | 1716 | slow |

### Analysis

Fastest benchmark **BenchmarkTransactionManagerGetTransaction-5** at 1108.00 ns/op (902527 ops/s). Slowest benchmark **BenchmarkNFTMarketplaceMint-5** at 2217720.00 ns/op (451 ops/s).
Average throughput 211570 ops/s, indicating moderate overall performance. The slowest benchmark suggests a ceiling of 2217720.00 ns/op (~451 ops/s).
The ns/op metric reflects the time each operation takes; lower ns/op and higher ops/s indicate better performance.

### Upgrade Plan

The following benchmarks are classified as slow and may need optimisation:

- BenchmarkTransactionManagerBurnAndRelease-5 (199231.00 ns/op)
- BenchmarkTransactionManagerLockAndMint-5 (200195.00 ns/op)
- BenchmarkTransactionHash-5 (1400781.00 ns/op)
- BenchmarkNFTMarketplaceMint-5 (2217720.00 ns/op)

### Bottleneck Analysis and Repair Plan

- BenchmarkTransactionManagerBurnAndRelease-5: reduce allocations and reuse objects; trim memory usage
- BenchmarkTransactionManagerLockAndMint-5: reduce allocations and reuse objects; trim memory usage
- BenchmarkTransactionHash-5: trim memory usage
- BenchmarkNFTMarketplaceMint-5: reduce allocations and reuse objects; trim memory usage
