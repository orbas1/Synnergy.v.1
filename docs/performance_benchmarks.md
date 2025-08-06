# Performance Benchmarks

This document captures baseline performance measurements for critical paths in
the Synnergy platform. The benchmarks provide a regression guard so that new
changes do not degrade throughput or memory usage. All benchmarks are executed
with Go's built‑in testing framework (`go test -bench`) and analysed with
[`benchstat`](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat).

## Running Benchmarks

From the repository root run:

```bash
go test -bench . -benchmem -run ^$
```

Individual packages or patterns can be supplied to `-bench`. Results should be
stored under the `benchmarks/` directory. For example, to update the transaction
manager baseline:

```bash
go test -bench=TransactionManager -benchmem -run ^$ ./core \
  | tee benchmarks/transaction_manager.txt
```

## Transaction Manager Baseline

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| LockAndMint | 4094 | 744 | 7 |
| BurnAndRelease | 3661 | 863 | 9 |
| ListTransactions | 857290 | 1286145 | 1 |
| GetTransaction | 447.5 | 0 | 0 |

### Performance Budget

The Transaction Manager operations should not regress more than **10%** from
these baseline values. Continuous Integration compares new benchmark runs
against `benchmarks/transaction_manager.txt` using `benchstat` and fails the
build if the threshold is exceeded.

## Adding New Benchmarks

1. Write a `BenchmarkXxx` function following the `testing` package
   conventions.
2. Run the benchmark and save the output under `benchmarks/<name>.txt`.
3. Update this document with a new table and describe any relevant performance
   budgets or expectations.

Keeping baselines current allows contributors to detect regressions early and
track the impact of optimisations over time.

