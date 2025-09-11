# Benchmark Full Report and Assessment

This document provides an enterprise-grade view of benchmark performance across the Synnergy codebase. Benchmarks are executed with `go test -bench . -benchmem ./...`.

Run `cd 'Security assessments & Benchmark assessments' && go run ./Benchmarks/cmd/runbench` to refresh results. The command rewrites this file with the latest measurements.

The `ns/op` column reports average time per operation, while `ops/s` derives throughput. Lower `ns/op` and higher `ops/s` indicate better performance.

## Automated Benchmark Results

| Benchmark | ns/op | ops/s | B/op | allocs/op | Rating |
|-----------|------|-------|------|-----------|--------|
| BenchmarkTransactionManagerGetTransaction-5 | 1611.00 | 620732 | 0 | 0 | fast |
| BenchmarkBaseTokenTransfer-5 | 2588.00 | 386399 | 0 | 0 | moderate |
| BenchmarkRegistryInfo-5 | 3154.00 | 317058 | 0 | 0 | moderate |
| BenchmarkTransactionManagerListTransactions-5 | 4684.00 | 213493 | 128 | 1 | moderate |
| BenchmarkLedgerApplyTransaction-5 | 6692.00 | 149432 | 80 | 6 | moderate |
| BenchmarkTransactionManagerLockAndMint-5 | 212099.00 | 4715 | 2600 | 17 | slow |
| BenchmarkTransactionHash-5 | 371829.00 | 2689 | 1128 | 10 | slow |
| BenchmarkTransactionManagerBurnAndRelease-5 | 374925.00 | 2667 | 2680 | 19 | slow |
| BenchmarkNFTMarketplaceMint-5 | 2803152.00 | 357 | 164920 | 1716 | slow |

### Analysis

Fastest benchmark **BenchmarkTransactionManagerGetTransaction-5** at 1611.00 ns/op (620732 ops/s). Slowest benchmark **BenchmarkNFTMarketplaceMint-5** at 2803152.00 ns/op (357 ops/s).
Average throughput 188616 ops/s, indicating moderate overall performance. The slowest benchmark suggests a ceiling of 2803152.00 ns/op (~357 ops/s).
The ns/op metric reflects the time each operation takes; lower ns/op and higher ops/s indicate better performance.

### Upgrade Plan

The following benchmarks are classified as slow and may need optimisation:

- BenchmarkTransactionManagerLockAndMint-5 (212099.00 ns/op)
- BenchmarkTransactionHash-5 (371829.00 ns/op)
- BenchmarkTransactionManagerBurnAndRelease-5 (374925.00 ns/op)
- BenchmarkNFTMarketplaceMint-5 (2803152.00 ns/op)

### Bottleneck Analysis and Repair Plan

- BenchmarkTransactionManagerLockAndMint-5: reduce allocations and reuse objects; trim memory usage
- BenchmarkTransactionHash-5: trim memory usage
- BenchmarkTransactionManagerBurnAndRelease-5: reduce allocations and reuse objects; trim memory usage
- BenchmarkNFTMarketplaceMint-5: reduce allocations and reuse objects; trim memory usage
