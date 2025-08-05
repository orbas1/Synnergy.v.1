
.PHONY: staticcheck gosec govulncheck security

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
