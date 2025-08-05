.PHONY: staticcheck gosec govulncheck security bench

staticcheck:
	staticcheck ./...

gosec:
	gosec ./...

govulncheck:
	govulncheck ./...

security: staticcheck gosec govulncheck

bench:
	go test -bench=TransactionManager -benchmem -run ^$ . | tee benchmarks/current.txt
