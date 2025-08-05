.PHONY: staticcheck gosec govulncheck security build build-experimental build-dev build-test build-prod

staticcheck:
	staticcheck ./...

gosec:
	gosec ./...

govulncheck:
	govulncheck ./...

security: staticcheck gosec govulncheck



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


