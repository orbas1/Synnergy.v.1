# Benchmark Full Report and Assessment

This document provides an enterprise-grade view of benchmark performance across the Synnergy codebase. Benchmarks are executed with `go test -bench . -benchmem ./...`.

Run `cd 'Security assessments & Benchmark assessments' && go run ./Benchmarks/cmd/runbench` to refresh results. The command rewrites this file with the latest measurements.

The `ns/op` column reports average time per operation, while `ops/s` derives throughput. Lower `ns/op` and higher `ops/s` indicate better performance.

## Automated Benchmark Results

| Benchmark | ns/op | ops/s | B/op | allocs/op | Rating |
|-----------|------|-------|------|-----------|--------|
| BenchmarkTransactionManagerGetTransaction-5 | 1323.00 | 755858 | 0 | 0 | fast |
| BenchmarkBaseTokenTransfer-5 | 2575.00 | 388350 | 0 | 0 | moderate |
| BenchmarkRegistryInfo-5 | 2962.00 | 337610 | 0 | 0 | moderate |
| BenchmarkTransactionManagerListTransactions-5 | 4777.00 | 209336 | 128 | 1 | moderate |
| BenchmarkLedgerApplyTransaction-5 | 4854.00 | 206016 | 80 | 6 | moderate |
| BenchmarkTransactionManagerBurnAndRelease-5 | 113135.00 | 8839 | 2680 | 19 | slow |
| BenchmarkTransactionManagerLockAndMint-5 | 210660.00 | 4747 | 2600 | 17 | slow |
| BenchmarkTransactionHash-5 | 1099397.00 | 910 | 1128 | 10 | slow |
| BenchmarkNFTMarketplaceMint-5 | 2801753.00 | 357 | 164448 | 1706 | slow |

### Analysis

Fastest benchmark **BenchmarkTransactionManagerGetTransaction-5** at 1323.00 ns/op (755858 ops/s). Slowest benchmark **BenchmarkNFTMarketplaceMint-5** at 2801753.00 ns/op (357 ops/s).
Average throughput 212447 ops/s, indicating moderate overall performance. The slowest benchmark suggests a ceiling of 2801753.00 ns/op (~357 ops/s).
The ns/op metric reflects the time each operation takes; lower ns/op and higher ops/s indicate better performance.

### Upgrade Plan

The following benchmarks are classified as slow and may need optimisation:

- BenchmarkTransactionManagerBurnAndRelease-5 (113135.00 ns/op)
- BenchmarkTransactionManagerLockAndMint-5 (210660.00 ns/op)
- BenchmarkTransactionHash-5 (1099397.00 ns/op)
- BenchmarkNFTMarketplaceMint-5 (2801753.00 ns/op)

### Bottleneck Analysis and Repair Plan

- BenchmarkTransactionManagerBurnAndRelease-5: reduce allocations and reuse objects; trim memory usage
- BenchmarkTransactionManagerLockAndMint-5: reduce allocations and reuse objects; trim memory usage
- BenchmarkTransactionHash-5: trim memory usage
- BenchmarkNFTMarketplaceMint-5: reduce allocations and reuse objects; trim memory usage
