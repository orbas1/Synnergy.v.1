
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



.PHONY: docs docs-serve

docs:
	mkdocs build

docs-serve:
	mkdocs serve

build:
	go build ./...

build-experimental:
	go build -tags experimental ./...

build-dev:
	go build -tags dev ./...

build-test:
	go build -tags test ./...

build-prod:
	go build -tags prod ./...



