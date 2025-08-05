# Performance Benchmarks

This document captures baseline performance measurements for critical paths in the Synnergy platform. Benchmarks are executed using Go's built-in benchmarking framework (`go test -bench`).

## Transaction Manager Benchmarks

```
go test -bench=TransactionManager -benchmem -run ^$ .
```

| Benchmark | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
| LockAndMint | 4094 | 744 | 7 |
| BurnAndRelease | 3661 | 863 | 9 |
| ListTransactions | 857290 | 1286145 | 1 |
| GetTransaction | 447.5 | 0 | 0 |

### Performance Budget

The Transaction Manager operations should not regress more than **10%** from these baseline values. Continuous Integration will compare new benchmark results against `benchmarks/transaction_manager.txt` using `benchstat`.

