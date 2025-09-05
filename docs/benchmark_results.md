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
| BenchmarkNFTMarketplaceMint | 1000000 | 1564 |

## synnergy/internal/tokens package

| Benchmark | Iterations | ns/op |
| --- | --- | --- |
| BenchmarkBaseTokenTransfer | 30524721 | 39.35 |
| BenchmarkRegistryInfo | 11217840 | 107.9 |

Benchmarks executed on `Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz`.

